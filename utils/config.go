package utils

import (
	databaseconfig "github.com/MonikaCat/ibcjuno/db/config"
	tokensconfig "github.com/MonikaCat/ibcjuno/token"
)

var (
	// Cfg represents the configuration used during the execution
	Cfg Config
)

// Config defines all necessary IBCJuno configuration parameters.
type Config struct {
	bytes []byte

	Database databaseconfig.Config     `yaml:"database"`
	Tokens   tokensconfig.TokensConfig `yaml:"tokens"`
}

// NewConfig builds a new Config instance
func NewConfig(
	dbConfig databaseconfig.Config,
	tokensConfig tokensconfig.TokensConfig,
) Config {
	return Config{
		Database: dbConfig,
		Tokens:   tokensConfig,
	}
}

func DefaultConfig() Config {
	return NewConfig(
		databaseconfig.DefaultDatabaseConfig(),
		tokensconfig.DefaultTokensConfig(),
	)
}

func (c Config) GetBytes() []byte {
	return c.bytes
}
