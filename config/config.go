package config

import (
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Config defines all necessary configuration parameters.
type Config struct {
	DB     DatabaseConfig `yaml:"database"`
	Tokens TokensConfig   `yaml:"tokens"`
}

// DatabaseConfig defines database connection parameters.

type DatabaseConfig struct {
	Name               string `yaml:"name"`
	Host               string `yaml:"host"`
	Port               int64  `yaml:"port"`
	User               string `yaml:"user"`
	Password           string `yaml:"password"`
	SSLMode            string `yaml:"ssl_mode,omitempty"`
	Schema             string `yaml:"schema,omitempty"`
	MaxOpenConnections int    `yaml:"max_open_connections"`
	MaxIdleConnections int    `yaml:"max_idle_connections"`
}

type TokensConfig struct {
	Tokens []Token `yaml:"token"`
}

type Token struct {
	Name  string      `yaml:"name"`
	Units []TokenUnit `yaml:"units"`
}

type TokenUnit struct {
	Denom    string   `yaml:"denom"`
	Exponent int      `yaml:"exponent"`
	Aliases  []string `yaml:"aliases"`
	PriceID  string   `yaml:"price_id"`
}

// ParseConfig reads and parse config.yaml file.
// An error is returned if the operation fails.
func ParseConfig(configPath string) Config {
	if configPath == "" {
		log.Fatal("invalid configuration file")
	}

	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to read config"))
	}

	var cfg Config
	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to unmarshal config"))
	}

	return cfg
}
