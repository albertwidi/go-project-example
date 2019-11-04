package envfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	yaml "gopkg.in/yaml.v2"
)

// EnvConfigYAML struct
type EnvConfigYAML struct {
	Envs []EnvYAML `yaml:"envs"`
}

// EnvYAML file struct
type EnvYAML struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// Load envfile
// capable of reading multiple files
func Load(files ...string) error {
	for _, envFile := range files {
		if envFile == "" {
			continue
		}

		ext := filepath.Ext(envFile)
		var kv map[string]string
		var err error

		switch ext {
		case ".toml":
			kv, err = loadToml(envFile)
		case ".yaml", ".yml":
			kv, err = loadYaml(envFile)
		default:
			err = fmt.Errorf("cannot process file with format %s", ext)
		}

		if err != nil {
			return err
		}

		// insert all value in the yaml file into env variable
		for k, v := range kv {
			if err := os.Setenv(strings.ToUpper(k), v); err != nil {
				return err
			}
		}
	}
	return nil
}

func loadToml(file string) (map[string]string, error) {
	kv := make(map[string]string)
	_, err := toml.DecodeFile(file, &kv)

	return kv, err
}

func loadYaml(file string) (map[string]string, error) {
	envs := EnvConfigYAML{}
	out, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(out, &envs); err != nil {
		return nil, err
	}

	kv := make(map[string]string)

	for _, e := range envs.Envs {
		kv[e.Name] = e.Value
	}

	return kv, nil
}
