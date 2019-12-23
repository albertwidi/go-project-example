package s3

import (
	"context"
	"errors"
	"os"

	"github.com/albertwidi/go-project-example/internal/pkg/objectstorage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"gocloud.dev/blob"
	"gocloud.dev/blob/s3blob"
)

// S3 struct
type S3 struct {
	storageBlobBucket *blob.Bucket
	config            *Config
}

// Config of S3
type Config struct {
	credentials *credentials.Credentials
	region      string
	endpoint    string
	bucket      string
	disableSSL  bool
	// Set this to `true` to force the request to use path-style addressing,
	// i.e., `http://s3.amazonaws.com/BUCKET/KEY`. By default, the S3 client
	// will use virtual hosted bucket addressing when possible
	// (`http://BUCKET.s3.amazonaws.com/KEY`).
	//
	// Note: This configuration option is specific to the Amazon S3 service.
	//
	// See http://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html
	// for Amazon S3: Virtual Hosting of Buckets
	forcePathStyle bool
	bucketProto    string
	bucketURL      string
}

// NewConfig return new configuration for s3
func NewConfig(ctx context.Context, creds *credentials.Credentials) (*Config, error) {
	c := Config{
		credentials: creds,
	}
	return &c, nil
}

// SetRegion to set s3 region
func (c *Config) SetRegion(region string) *Config {
	c.region = region
	return c
}

// SetEndpoint to set s3 endpoint
func (c *Config) SetEndpoint(endpoint string) *Config {
	c.endpoint = endpoint
	return c
}

// DisableSSL to disable SSL when connecting to s3
func (c *Config) DisableSSL(disable bool) *Config {
	c.disableSSL = disable
	return c
}

// ForcePathStyle when get/put object from s3
func (c *Config) ForcePathStyle(force bool) *Config {
	c.forcePathStyle = force
	return c
}

// SetBucket for storing s3 bucket information
func (c *Config) SetBucket(bucket string) *Config {
	c.bucket = bucket
	return c
}

// SetBucketProto to set the proto when connecting to s3 bucket
func (c *Config) SetBucketProto(proto string) *Config {
	c.bucketProto = proto
	return c
}

// SetBucketURL for s3
func (c *Config) SetBucketURL(url string) *Config {
	c.bucketURL = url
	return c
}

// New S3 storage
func New(ctx context.Context, config *Config) (*S3, error) {
	if config == nil {
		return nil, errors.New("s3 config cannot be nil")
	}
	c := aws.Config{
		Region:           aws.String(config.region),
		Credentials:      config.credentials,
		DisableSSL:       aws.Bool(config.disableSSL),
		S3ForcePathStyle: aws.Bool(config.forcePathStyle),
	}

	// if we want to use digitalocean or another api that compatible to s3
	if config.endpoint != "" {
		c.Endpoint = aws.String(config.endpoint)
	}
	sess, err := session.NewSession(&c)
	if err != nil {
		return nil, err
	}
	bb, err := s3blob.OpenBucket(ctx, sess, config.bucket, nil)
	if err != nil {
		return nil, err
	}

	s := S3{
		storageBlobBucket: bb,
		config:            config,
	}
	return &s, nil
}

// Bucket function
func (s3 *S3) Bucket() *blob.Bucket {
	return s3.storageBlobBucket
}

// Name of objectstorage provider
func (s3 *S3) Name() string {
	return objectstorage.StorageS3
}

// BucketName return name of bucket used for objectstorage
func (s3 *S3) BucketName() string {
	return s3.config.bucket
}

// BucketURL return full path for bucket
func (s3 *S3) BucketURL() string {
	return ""
	// return fmt.Sprintf("%s%s.%s", s3.config.bucketProto, s3.BucketName(), s3.config.bucketURL)
}

// Close will close the s3 bucket and return error
func (s3 *S3) Close() error {
	return s3.Bucket().Close()
}

// CredentialsFromClient is a helper function with clientID and clientSecret provided from client
func CredentialsFromClient(ctx context.Context, clientID, clientSecret, token string) (*credentials.Credentials, error) {
	creds := credentials.NewStaticCredentials(clientID, clientSecret, token)
	return creds, nil
}

// CredentialsFromSharedProfile is a helper function to load credentials from a profile file
// usually the profile file located at $HOME/.aws/profile
func CredentialsFromSharedProfile(ctx context.Context, filename, profile string) (*credentials.Credentials, error) {
	// check whether the file is exists or not
	_, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	creds := credentials.NewSharedCredentials(filename, profile)
	return creds, nil
}
