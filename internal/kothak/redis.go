package kothak

// Redis interface for infra
type Redis interface {
}

// RedisConfig of kothak
type RedisConfig struct {
	MaxIdle   int               `json:"max_idle_conn" yaml:"max_idle_conn" toml:"max_idle_conn"`
	MaxActive int               `json:"max_active_conn" yaml:"max_active_conn" toml:"max_active_conn"`
	Timeout   int               `json:"timeout" yaml:"timeout" toml:"timeout"`
	Rds       []RedisConnConfig `json:"connect" yaml:"connect" toml:"connect"`
}

// RedisConnConfig struct
type RedisConnConfig struct {
	Name      string `json:"name" yaml:"name" toml:"name"`
	Address   string `json:"address" yaml:"address" toml:"address"`
	MaxIdle   int    `json:"max_idle_conn" yaml:"max_idle_conn" toml:"max_idle_conn"`
	MaxActive int    `json:"max_active_conn" yaml:"max_active_conn" toml:"max_active_conn"`
	Timeout   int    `json:"timeout" yaml:"timeout" toml:"timeout"`
}
