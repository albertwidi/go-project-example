// expect this program to run from project root
// go run internal/config/test_generator/main.go

package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/albertwidi/go-project-example/internal/config"
)

func main() {
	// the path is not relative, run from project root
	configFile := "project.config.toml"
	envFile := "project.env.toml"

	defaultConfig := config.DefaultConfig{}
	if err := config.ParseFile(configFile, &defaultConfig, envFile); err != nil {
		log.Fatal(err)
	}

	// json
	out, err := json.MarshalIndent(defaultConfig, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("config.test.json")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := f.Write(out); err != nil {
		log.Fatal(err)
	}
	// yaml
	// toml
}
