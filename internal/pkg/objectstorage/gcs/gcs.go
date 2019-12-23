package gcs

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/albertwidi/go-project-example/internal/pkg/objectstorage"
	"gocloud.dev/blob"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/gcp"
	"golang.org/x/oauth2/google"
)

// error list
var (
	ErrCredentialsEmpty = errors.New("credentials is empty")
)

// Config of Google Cloud Storage
type Config struct {
	Bucket      string
	BucketProto string
	BucketURL   string
	credentials *google.Credentials
}

// NewConfig to build a new config
func NewConfig(ctx context.Context, creds *google.Credentials) (*Config, error) {
	var (
		err error
	)

	if creds == nil {
		// fallback to credentials from gcloud if empty
		creds, err = CredentialsFromGcloud(ctx)
		if err != nil {
			return nil, err
		}
	}

	c := Config{
		Bucket:      "",
		BucketProto: "gs://",
		BucketURL:   "gcr.io",
		credentials: creds,
	}

	return &c, nil
}

// SetBucket set bucket name
func (c *Config) SetBucket(bucket string) *Config {
	c.Bucket = bucket
	return c
}

// SetBucketProto set bucket protocol
func (c *Config) SetBucketProto(bucketProto string) *Config {
	c.BucketProto = bucketProto
	return c
}

// SetBucketURL set bucket URL
func (c *Config) SetBucketURL(bucketURL string) *Config {
	c.BucketURL = bucketURL
	return c
}

// GCS struct
type GCS struct {
	config            *Config
	baseURL           string
	storageBlobBucket *blob.Bucket
	credentials       *google.Credentials
	httpClient        *gcp.HTTPClient
}

// New gcs storage
// gcs storage expect user already authenticated using 'gcs' command
// json path is a path to json service account path
func New(ctx context.Context, config *Config) (*GCS, error) {
	if config == nil {
		return nil, errors.New("gcs config cannot be nil")
	}

	var (
		creds *google.Credentials
		err   error
	)

	creds = config.credentials
	if creds == nil {
		return nil, errors.New("gcloud credentials is nil")
	}

	client, err := gcp.NewHTTPClient(gcp.DefaultTransport(), gcp.CredentialsTokenSource(creds))
	if err != nil {
		return nil, err
	}

	bb, err := gcsblob.OpenBucket(ctx, client, config.Bucket, nil)
	if err != nil {
		return nil, err
	}

	gcs := GCS{
		config:            config,
		storageBlobBucket: bb,
		httpClient:        client,
		baseURL:           fmt.Sprintf("%s%s", config.BucketProto, config.Bucket),
	}
	return &gcs, nil
}

// Bucket function
func (gcs *GCS) Bucket() *blob.Bucket {
	return gcs.storageBlobBucket
}

// Name of objectstorage provider
func (gcs *GCS) Name() string {
	return objectstorage.StorageGCS
}

// BucketName return name of bucket used for objectstorage
func (gcs *GCS) BucketName() string {
	return gcs.config.Bucket
}

// BucketURL return full path for bucket
func (gcs *GCS) BucketURL() string {
	return fmt.Sprintf("%s.%s", gcs.baseURL, gcs.config.BucketURL)
}

// Close the gcs bucket
func (gcs *GCS) Close() error {
	return gcs.Bucket().Close()
}

// CredentialsFromFile return content of credentials from file
func CredentialsFromFile(ctx context.Context, path string) (*google.Credentials, error) {
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(jsonData) == 0 {
		return nil, ErrCredentialsEmpty
	}

	creds, err := google.CredentialsFromJSON(ctx, []byte(jsonData), "https://www.googleapis.com/auth/cloud-platform")
	return creds, err
}

// CredentialsFromEnvVar return content of credentials from environment variable
func CredentialsFromEnvVar(ctx context.Context, name string) (*google.Credentials, error) {
	content := os.Getenv(name)
	if len(content) == 0 {
		return nil, ErrCredentialsEmpty
	}

	creds, err := google.CredentialsFromJSON(ctx, []byte(content), "https://www.googleapis.com/auth/cloud-platform")
	return creds, err
}

// CredentialsFromGcloud to provide credentials form default gcloud command
func CredentialsFromGcloud(ctx context.Context) (*google.Credentials, error) {
	creds, err := gcp.DefaultCredentials(ctx)
	return creds, err
}

// CredentialsFromString if user want to load credentials by their own
func CredentialsFromString(ctx context.Context, content string) (*google.Credentials, error) {
	creds, err := google.CredentialsFromJSON(ctx, []byte(content), "https://www.googleapis.com/auth/cloud-platform")
	return creds, err
}
