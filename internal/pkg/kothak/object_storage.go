package kothak

// ObjectStorageConfig struct
type ObjectStorageConfig struct {
	Name           string    `yaml:"name"`
	Provider       string    `yaml:"provider"`
	Region         string    `yaml:"region"`
	Endpoint       string    `yaml:"endpoint"`
	Bucket         string    `yaml:"bucket"`
	ClientID       string    `yaml:"client_id"`
	ClientSecret   string    `yaml:"client_secret"`
	JSONKey        string    `yaml:"json_key"`
	DisableSSL     bool      `yaml:"disable_ssl"`
	ForcePathStyle bool      `yaml:"force_path_style"`
	BucketProto    string    `yaml:"bucket_proto"`
	BucketURL      string    `yaml:"bucket_url"`
	S3             S3Config  `yaml:"s3"`
	GCS            GCSConfig `yaml:"gcs"`
}

// S3Config for s3 storage
type S3Config struct {
	ClientID       string `yaml:"client_id"`
	ClientSecret   string `yaml:"client_secret"`
	DisableSSL     bool   `yaml:"disable_ssl"`
	ForcePathStyle bool   `yaml:"force_path_style"`
}

// GCSConfig for google cloud storage
type GCSConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	JSONKey      string `yaml:"json_key"`
}
