package kothak

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/albertwidi/go-project-example/internal/pkg/log/logger"
	"github.com/albertwidi/go-project-example/internal/pkg/objectstorage"
	"github.com/albertwidi/go-project-example/internal/pkg/objectstorage/gcs"
	"github.com/albertwidi/go-project-example/internal/pkg/objectstorage/local"
	"github.com/albertwidi/go-project-example/internal/pkg/objectstorage/s3"
	"github.com/albertwidi/go-project-example/internal/pkg/redis"
	redigo "github.com/albertwidi/go-project-example/internal/pkg/redis/redigo"
	"github.com/albertwidi/go-project-example/internal/pkg/sqldb"
	"github.com/jmoiron/sqlx"
	"go.opencensus.io/trace"
)

// Config of kothak
type Config struct {
	DBConfig            DBConfig              `json:"database" yaml:"database" toml:"database"`
	RedisConfig         RedisConfig           `json:"redis" yaml:"redis" toml:"redis"`
	ObjectStorageConfig []ObjectStorageConfig `json:"object_storage" yaml:"object_storage" toml:"object_storage"`
}

// Kothak struct
type Kothak struct {
	objStorages map[string]*objectstorage.Storage
	dbs         map[string]*sqldb.DB
	rds         map[string]redis.Redis
	logger      logger.Logger
	// this mutex is used in two place
	// the usage shared because the usage is not collide
	// 1. when we initialize all the connections
	// 2. when we want to get the connection
	mutex sync.Mutex
}

func (k *Kothak) setSQLDB(name string, db *sqldb.DB) {
	k.mutex.Lock()
	k.dbs[name] = db
	k.mutex.Unlock()
}

func (k *Kothak) setRedis(name string, rds redis.Redis) {
	k.mutex.Lock()
	k.rds[name] = rds
	k.mutex.Unlock()
}

func (k *Kothak) setObjectStorage(name string, obj objectstorage.StorageProvider) {
	k.mutex.Lock()
	k.objStorages[name] = objectstorage.New(obj)
	k.mutex.Unlock()
}

// New kothak instance
func New(ctx context.Context, kothakConfig Config, logger logger.Logger) (*Kothak, error) {
	ctx, span := trace.StartSpan(ctx, "ktohak/new")
	defer span.End()

	var (
		kothak = Kothak{
			objStorages: make(map[string]*objectstorage.Storage),
			dbs:         make(map[string]*sqldb.DB),
			rds:         make(map[string]redis.Redis),
			logger:      logger,
		}

		group = sync.WaitGroup{}
		errs  []error
		err   error
	)

	// set default configuration for DBConfig
	if err := kothakConfig.DBConfig.SetDefault(); err != nil {
		return nil, err
	}

	// connect to object storage
	for _, objStorageConfig := range kothakConfig.ObjectStorageConfig {
		group.Add(1)
		go func(config ObjectStorageConfig) {
			_, span = trace.StartSpan(ctx, fmt.Sprintf("object_storage/init/%s", config.Name))
			defer func() {
				span.End()
				group.Done()
			}()

			var (
				provider objectstorage.StorageProvider
				err      error
			)

			switch strings.ToLower(config.Provider) {
			// local storage
			case objectstorage.StorageLocal:
				// defaulted to not delete local bucket when close the program
				provider, err = local.New(ctx, fmt.Sprintf("./%s", config.Bucket), &local.Options{DeleteOnClose: false})

			// gcs compatible storage
			case objectstorage.StorageGCS:
				gcsCreds, err := gcs.CredentialsFromFile(ctx, config.GCS.JSONKey)
				if err != nil {
					errs = append(errs, err)
					return
				}

				gcsConfig, err := gcs.NewConfig(ctx, gcsCreds)
				if err != nil {
					errs = append(errs, err)
					return
				}
				gcsConfig.
					SetBucket(config.Bucket).
					SetBucketProto(config.BucketProto).
					SetBucketURL(config.BucketURL)

				provider, err = gcs.New(ctx, gcsConfig)
				if err != nil {
					errs = append(errs, err)
					return
				}

			// s3 compatible storage
			case objectstorage.StorageS3, objectstorage.StorageDO, objectstorage.StorageMinio:
				s3Creds, err := s3.CredentialsFromClient(ctx, config.S3.ClientID, config.S3.ClientSecret, "")
				if err != nil {
					errs = append(errs, err)
					return
				}

				s3Config, err := s3.NewConfig(ctx, s3Creds)
				if err != nil {
					errs = append(errs, err)
					return
				}

				s3Config.
					SetBucket(config.Bucket).
					SetBucketProto(config.BucketProto).
					SetBucketURL(config.BucketURL).
					SetRegion(config.Region).
					SetEndpoint(config.Endpoint).
					DisableSSL(config.S3.DisableSSL).
					ForcePathStyle(config.S3.ForcePathStyle)

				provider, err = s3.New(ctx, s3Config)
				if err != nil {
					errs = append(errs, err)
					return
				}

			default:
				err = errors.New("kothak: object storage provider not found")
				errs = append(errs, err)
				return
			}

			if err != nil {
				errs = append(errs, err)
				return
			}

			logger.Debugf("kothak: Connected to object_storage %s", config.Name)

			kothak.setObjectStorage(config.Name, provider)
		}(objStorageConfig)
	}

	// connect to redis
	for _, redisconfig := range kothakConfig.RedisConfig.Rds {
		group.Add(1)
		go func(redisconfig RedisConnConfig) {
			_, span = trace.StartSpan(ctx, fmt.Sprintf("redis/init/%s", redisconfig.Name))
			defer func() {
				group.Done()
				span.End()
			}()

			conf := redigo.Config{
				MaxActive: kothakConfig.RedisConfig.MaxActive,
				MaxIdle:   kothakConfig.RedisConfig.MaxIdle,
				Timeout:   kothakConfig.RedisConfig.Timeout,
			}

			r, err := redigo.New(ctx, redisconfig.Address, &conf)
			if err != nil {
				errs = append(errs, err)
				return
			}

			logger.Debugf("Kothak: Connected to Redis %s", redisconfig.Name)

			kothak.setRedis(redisconfig.Name, r)
		}(redisconfig)
	}

	// connect to database
	for _, dbconfig := range kothakConfig.DBConfig.SQLDBs {
		group.Add(1)
		go func(dbconfig SQLDBConfig) {
			_, span = trace.StartSpan(ctx, fmt.Sprintf("database/connect/%s", dbconfig.Name))
			defer func() {
				group.Done()
				span.End()
			}()

			var (
				err        error
				leaderDB   *sqlx.DB
				followerDB *sqlx.DB
			)

			if dbconfig.Driver == "" {

			}

			// setup leader connection
			if err := dbconfig.LeaderConnConfig.SetDefault(kothakConfig.DBConfig); err != nil {
				errs = append(errs, err)
				return
			}
			// connect to leader
			leaderDB, err = sqldb.Connect(ctx, dbconfig.Driver, dbconfig.LeaderConnConfig.DSN, &sqldb.ConnectOptions{
				Retry:              dbconfig.LeaderConnConfig.MaxRetry,
				MaxOpenConnections: dbconfig.LeaderConnConfig.MaxOpenConnections,
				MaxIdleConnections: dbconfig.LeaderConnConfig.MaxIdleConnections,
			})
			if err != nil {
				errs = append(errs, err)
				return
			}
			// by default, set replica to leader
			followerDB = leaderDB

			// connect to replica
			if dbconfig.ReplicaConnConfig.DSN != "" {
				if err := dbconfig.ReplicaConnConfig.SetDefault(kothakConfig.DBConfig); err != nil {
					errs = append(errs, err)
					return
				}
				followerDB, err = sqldb.Connect(ctx, dbconfig.Driver, dbconfig.ReplicaConnConfig.DSN, &sqldb.ConnectOptions{
					Retry:              dbconfig.ReplicaConnConfig.MaxRetry,
					MaxOpenConnections: dbconfig.ReplicaConnConfig.MaxOpenConnections,
					MaxIdleConnections: dbconfig.ReplicaConnConfig.MaxIdleConnections,
				})
				if err != nil {
					errs = append(errs, err)
					return
				}
			}

			db, err := sqldb.Wrap(ctx, leaderDB, followerDB)
			if err != nil {
				errs = append(errs, err)
				return
			}

			logger.Debugf("kothak: connected to DB %s", dbconfig.Name)

			kothak.setSQLDB(dbconfig.Name, db)
		}(dbconfig)
	}

	// wait for all connections
	group.Wait()
	// check for error, if error length is greater than 1
	// set err to errs[0]
	if len(errs) > 0 {
		err = errs[0]
	}

	return &kothak, err
}

// CloseAll to close all connected resources
// TODO: check error when closing connections and close connection concurrently
func (k *Kothak) CloseAll() error {
	for _, objStorage := range k.objStorages {
		objStorage.Close()
	}

	for _, sqldb := range k.dbs {
		sqldb.Close()
	}

	for _, redis := range k.rds {
		redis.Close()
	}
	return nil
}

// GetSQLDB from kothak object
func (k *Kothak) GetSQLDB(dbname string) (*sqldb.DB, error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	i, ok := k.dbs[dbname]
	if !ok {
		err := fmt.Errorf("kothak: sql database with name %s does not exists", dbname)
		return nil, err
	}
	return i, nil
}

// MustGetSQLDB from kothak object
func (k *Kothak) MustGetSQLDB(dbname string) *sqldb.DB {
	db, err := k.GetSQLDB(dbname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return db
}

// GetRedis from kothak object
func (k *Kothak) GetRedis(redisname string) (redis.Redis, error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	i, ok := k.rds[redisname]
	if !ok {
		err := fmt.Errorf("kothak: redis with name %s does not exists", redisname)
		return nil, err
	}
	return i, nil
}

// MustGetRedis from kothak object
func (k *Kothak) MustGetRedis(redisname string) redis.Redis {
	r, err := k.GetRedis(redisname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return r
}

// GetObjectStorage from kothak object
func (k *Kothak) GetObjectStorage(objStorageName string) (*objectstorage.Storage, error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	i, ok := k.objStorages[objStorageName]
	if !ok {
		err := fmt.Errorf("kothak: object storage with name %s does not exists", objStorageName)
		return nil, err
	}
	return i, nil
}

// MustGetObjectStorage from kothak object
func (k *Kothak) MustGetObjectStorage(objStorageName string) *objectstorage.Storage {
	o, err := k.GetObjectStorage(objStorageName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return o
}
