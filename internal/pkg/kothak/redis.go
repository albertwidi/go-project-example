package kothak

// Redis interface for infra
type Redis interface {
}

// RedisConfig of kothak
type RedisConfig struct {
	MaxIdle   int               `yaml:"max_idle_conn"`
	MaxActive int               `yaml:"max_active_conn"`
	Timeout   int               `yaml:"timeout"`
	Rds       []RedisConnConfig `yaml:"connect"`
}

// RedisConnConfig struct
type RedisConnConfig struct {
	Name      string `yaml:"name"`
	Address   string `yaml:"address"`
	MaxIdle   int    `yaml:"max_idle_conn"`
	MaxActive int    `yaml:"max_active_conn"`
	Timeout   int    `yaml:"timeout"`
}
