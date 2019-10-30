package kothak

// Redis interface for infra
type Redis interface {
}

// RedisConfig of kothak
type RedisConfig struct {
	MaxIdle   int               `yaml:"max_idle_conn" toml:"max_idle"`
	MaxActive int               `yaml:"max_active_conn" toml:"max_active"`
	Timeout   int               `yaml:"timeout" toml:"timeout"`
	Rds       []RedisConnConfig `yaml:"connect" toml:"connect"`
}

// RedisConnConfig struct
type RedisConnConfig struct {
	Name      string `yaml:"name" toml:"name"`
	Address   string `yaml:"address" toml:"address"`
	MaxIdle   int    `yaml:"max_idle_conn" toml:"max_idle"`
	MaxActive int    `yaml:"max_active_conn" toml:"max_active"`
	Timeout   int    `yaml:"timeout" toml:"timeout"`
}
