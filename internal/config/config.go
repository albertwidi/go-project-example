package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/albertwidi/go-project-example/internal/kothak"
	"github.com/albertwidi/go-project-example/internal/pkg/envfile"
	"github.com/albertwidi/go-project-example/internal/pkg/tempe"
	"gopkg.in/yaml.v2"
)

// DefaultConfig for the project
type DefaultConfig struct {
	Servers   DefaultServers `json:"servers" yaml:"servers" toml:"servers"`
	Log       DefaultLog     `json:"log" yaml:"log" toml:"log"`
	Resources kothak.Config  `json:"resources" yaml:"resources" toml:"resources"`
}

// DefaultLog config for the project
type DefaultLog struct {
	Level string `json:"level" yaml:"level" toml:"level"`
	File  string `json:"file" yaml:"file" toml:"file"`
	Color bool   `json:"use_color" yaml:"use_color" toml:"use_color"`
}

// DefaultServers struct
type DefaultServers struct {
	Main  ServerConfig `json:"main" yaml:"main" toml:"main"`
	Debug ServerConfig `json:"debug" yaml:"debug" toml:"debug"`
	Admin ServerConfig `json:"admin" yaml:"admin" toml:"admin"`
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

// Print configuration in json schema
// all config with tag=protected:1 will be hidden from the print
func Print(v interface{}) {

}
