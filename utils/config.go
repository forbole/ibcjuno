package utils

import (
	apiconfig "github.com/forbole/ibcjuno/api/config"
	databaseconfig "github.com/forbole/ibcjuno/database/config"
)

var (
	// Cfg represents the configuration used during the execution
	Cfg Config
)

// Config defines all necessary IBCJuno configuration parameters.
type Config struct {
	bytes []byte

	API      apiconfig.APIConfig           `yaml:"api"`
	Database databaseconfig.DatabaseConfig `yaml:"database"`
}

// NewConfig builds new Config instance
func NewConfig(
	api apiconfig.APIConfig,
	dbConfig databaseconfig.DatabaseConfig,
) Config {
	return Config{
		API:      api,
		Database: dbConfig,
	}
}

// DefaultConfig returns default Config instance
func DefaultConfig() Config {
	return NewConfig(
		apiconfig.DefaultAPIConfig(),
		databaseconfig.DefaultDatabaseConfig(),
	)
}

// GetBytes returns slice of byte
func (c Config) GetBytes() []byte {
	return c.bytes
}
