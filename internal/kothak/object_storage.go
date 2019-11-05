package kothak

// ObjectStorageConfig struct
type ObjectStorageConfig struct {
	Name        string    `yaml:"name" toml:"name"`
	Provider    string    `yaml:"provider" toml:"provider"`
	Region      string    `yaml:"region" toml:"region"`
	Endpoint    string    `yaml:"endpoint" toml:"endpoint"`
	Bucket      string    `yaml:"bucket" toml:"bucket"`
	BucketProto string    `yaml:"bucket_proto" toml:"bucket_proto"`
	BucketURL   string    `yaml:"bucket_url" toml:"bucket_url"`
	S3          S3Config  `yaml:"s3" toml:"s3"`
	GCS         GCSConfig `yaml:"gcs" toml:"gcs"`
}

// S3Config for s3 storage
type S3Config struct {
	ClientID       string `yaml:"client_id" toml:"client_id"`
	ClientSecret   string `yaml:"client_secret" toml:"client_secret"`
	DisableSSL     bool   `yaml:"disable_ssl" toml:"disable_ssl"`
	ForcePathStyle bool   `yaml:"force_path_style" toml:"force_path_style"`
}

// GCSConfig for google cloud storage
type GCSConfig struct {
	ClientID     string `yaml:"client_id" toml:"client_id"`
	ClientSecret string `yaml:"client_secret" toml:"client_secret"`
	JSONKey      string `yaml:"json_key" toml:"json_key"`
}
