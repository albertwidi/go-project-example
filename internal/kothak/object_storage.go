package kothak

// ObjectStorageConfig struct
type ObjectStorageConfig struct {
	Name        string    `json:"name" yaml:"name" toml:"name"`
	Provider    string    `json:"provider" yaml:"provider" toml:"provider"`
	Region      string    `json:"region" yaml:"region" toml:"region"`
	Endpoint    string    `json:"endpoint" yaml:"endpoint" toml:"endpoint"`
	Bucket      string    `json:"bucket" yaml:"bucket" toml:"bucket"`
	BucketProto string    `json:"bucket_proto" yaml:"bucket_proto" toml:"bucket_proto"`
	BucketURL   string    `json:"bucket_url" yaml:"bucket_url" toml:"bucket_url"`
	S3          S3Config  `json:"s3" yaml:"s3" toml:"s3"`
	GCS         GCSConfig `json:"gcs" yaml:"gcs" toml:"gcs"`
}

// S3Config for s3 storage
type S3Config struct {
	ClientID       string `json:"client_id" yaml:"client_id" toml:"client_id"`
	ClientSecret   string `json:"client_secret" yaml:"client_secret" toml:"client_secret"`
	DisableSSL     bool   `json:"disable_ssl" yaml:"disable_ssl" toml:"disable_ssl"`
	ForcePathStyle bool   `json:"force_path_style" yaml:"force_path_style" toml:"force_path_style"`
}

// GCSConfig for google cloud storage
type GCSConfig struct {
	ClientID     string `json:"client_id" yaml:"client_id" toml:"client_id"`
	ClientSecret string `json:"client_secret" yaml:"client_secret" toml:"client_secret"`
	JSONKey      string `json:"json_key" yaml:"json_key" toml:"json_key"`
}
