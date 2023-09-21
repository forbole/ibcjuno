package utils

import (
	apiconfig "github.com/forbole/ibcjuno/api/config"
	databaseconfig "github.com/forbole/ibcjuno/database/config"
	parserconfig "github.com/forbole/ibcjuno/parser/config"
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
	Parser   parserconfig.ParserConfig     `yaml:"parser"`
}

// NewConfig builds new Config instance
func NewConfig(
	api apiconfig.APIConfig,
	dbConfig databaseconfig.DatabaseConfig,
	parserConfig parserconfig.ParserConfig,
) Config {
	return Config{
		API:      api,
		Database: dbConfig,
		Parser:   parserConfig,
	}
}

// DefaultConfig returns default Config instance
func DefaultConfig() Config {
	return NewConfig(
		apiconfig.DefaultAPIConfig(),
		databaseconfig.DefaultDatabaseConfig(),
		parserconfig.DefaultParserConfig(),
	)
}

// GetBytes returns slice of byte
func (c Config) GetBytes() []byte {
	return c.bytes
}
