package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

var supportedDrivers = map[string]string{
	"postgres": "github.com/lib/pq",
}

type otterConfig struct {
	moduleName    string
	MigrationsDir string `json:"migrationsDir"`
	DbDriver      string `json:"dbDriver"`
}

var config otterConfig

func readConfig() otterConfig {
	empty := otterConfig{}
	if config != empty {
		return config
	}
	readModuleName()
	readJsonConfig()

	if len(config.MigrationsDir) == 0 {
		config.MigrationsDir = "migrations"
	}
	if config.MigrationsDir[0] == '/' {
		panic("migrations directory must be relative to the root project, can't start with `/`")
	}
	config.MigrationsDir = strings.TrimPrefix(config.MigrationsDir, "./")
	config.MigrationsDir = strings.TrimSuffix(config.MigrationsDir, "/")

	if len(config.DbDriver) != 0 {
		isSupported := false
		for driver, _ := range supportedDrivers {
			if driver == config.DbDriver {
				isSupported = true
			}
		}
		if !isSupported {
			panic(fmt.Sprintf("unsuported driver `%s`, supported drivers are [postgres]", config.DbDriver))
		}
	}

	return config

}

func readJsonConfig() {
	f, err := os.Open("./otter.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		panic(err)
	}
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		panic(err)
	}
}

func readModuleName() {
	if len(config.moduleName) > 0 {
		return
	}
	f, err := os.Open("./go.mod")
	if err != nil {
		panic(err)
	}

	b := bufio.NewReader(f)

	data, err := b.ReadString('\n')
	if err != nil {
		panic(err)
	}
	for i, c := range data {
		if i < 7 {
			continue
		}
		if c == '\n' {
			break
		}
		config.moduleName += string(c)
	}
}
