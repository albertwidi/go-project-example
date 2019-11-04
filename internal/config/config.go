package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/albertwidi/go_project_example/internal/pkg/envfile"
	"github.com/albertwidi/go_project_example/internal/pkg/kothak"
	"github.com/albertwidi/go_project_example/internal/pkg/tempe"
	"gopkg.in/yaml.v2"
)

// DefaultConfig for the project
type DefaultConfig struct {
	Servers   DefaultServers `yaml:"servers" toml:"servers"`
	Resources kothak.Config  `yaml:"resources" toml:"resources"`
}

// DefaultServers struct
type DefaultServers struct {
	Main  ServerConfig `yaml:"main" toml:"main"`
	Dev   ServerConfig `yaml:"dev" toml:"dev"`
	Admin ServerConfig `yaml:"admin" toml:"admin"`
}

// ServerConfig struct
type ServerConfig struct {
	Address string `yaml:"address" toml:"address"`
}

// ParseFile for parsing config file and return DefaultConfig struct
func ParseFile(configFile string, dest interface{}, envFiles ...string) error {
	// prepare to replace ${ENV_VAR_NAME} with environment variable
	t, err := tempe.New(tempe.EnvVarPattern, tempe.EnvVarReplacerFunc)
	if err != nil {
		return err
	}

	if err := envfile.Load(envFiles...); err != nil {
		return err
	}

	out, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	// replacing with environment variables
	out, err = t.ReplaceBytes(out)
	if err != nil {
		return err
	}

	log.Printf("%s", string(out))

	ext := filepath.Ext(configFile)
	switch ext {
	case ".toml":
		err = toml.Unmarshal(out, dest)
	case ".yaml":
		err = yaml.Unmarshal(out, dest)
	default:
		err = fmt.Errorf("config: file ext is not valid. got %s", ext)
	}

	if err != nil {
		return err
	}

	return nil
}
