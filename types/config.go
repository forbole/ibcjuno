package config

import (
	databaseconfig "github.com/MonikaCat/ibc-token/database/config"
	// nodeconfig "github.com/forbole/juno/v3/node/config"
)

var (
	// Cfg represents the configuration used during the execution
	Cfg Config
)

// Config defines all necessary juno configuration parameters.
type Config struct {
	bytes []byte

	// Node     nodeconfig.Config     `yaml:"node"`
	Database databaseconfig.Config `yaml:"database"`
}

// NewConfig builds a new Config instance
func NewConfig(
	// nodeCfg nodeconfig.Config, 
	dbConfig databaseconfig.Config,
) Config {
	return Config{
		// Node:     nodeCfg,
		Database: dbConfig,
	}
}

func DefaultConfig() Config {
	return NewConfig(
		// nodeconfig.DefaultConfig(), 
		databaseconfig.DefaultDatabaseConfig(),
	)
}

func (c Config) GetBytes() []byte {
	return c.bytes
}
