package config

import (
	"github.com/albertwidi/go_project_example/internal/pkg/kothak"
)

// DefaultConfig for the project
type DefaultConfig struct {
	Servers   DefaultServers `toml:"servers"`
	Resources kothak.Config  `toml:"resources"`
}

// DefaultServers struct
type DefaultServers struct {
	Main  ServerConfig `toml:"main"`
	Dev   ServerConfig `toml:"dev"`
	Admin ServerConfig `toml:"admin"`
}

// ServerConfig struct
type ServerConfig struct {
	Address string `toml:"address"`
}
