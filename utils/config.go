package utils

import (
	databaseconfig "github.com/MonikaCat/ibcjuno/db/config"
	tokensconfig "github.com/MonikaCat/ibcjuno/types"
)

var (
	// Cfg represents the configuration used during the execution
	Cfg Config
)

// Config defines all necessary IBCJuno configuration parameters.
type Config struct {
	bytes []byte

	Database databaseconfig.DatabaseConfig `yaml:"database"`
	Tokens   tokensconfig.TokensConfig     `yaml:"tokens"`
}

// NewConfig builds new Config instance
func NewConfig(
	dbConfig databaseconfig.DatabaseConfig,
	tokensConfig tokensconfig.TokensConfig,
) Config {
	return Config{
		Database: dbConfig,
		Tokens:   tokensConfig,
	}
}

// DefaultConfig returns default Config instance
func DefaultConfig() Config {
	return NewConfig(
		databaseconfig.DefaultDatabaseConfig(),
		tokensconfig.DefaultTokensConfig(),
	)
}

// GetBytes returns slice of byte
func (c Config) GetBytes() []byte {
	return c.bytes
}
