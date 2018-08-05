package server

var Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Master   string `yaml:"master"`
	Follower string `yaml:"follower"`
}

type RedisConfig struct {
	Address string `yaml:"address"`
}
