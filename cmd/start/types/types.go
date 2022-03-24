package types

import (
	database "github.com/MonikaCat/ibcjuno/db"
	builder "github.com/MonikaCat/ibcjuno/db/builder"
	utils "github.com/MonikaCat/ibcjuno/utils"
)

// Config contains all the configuration for the "start" command
type Config struct {
	configParser utils.Parser
	buildDb      database.Builder
}

// NewConfig allows to build a new Config instance
func NewConfig() *Config {
	return &Config{}
}

// WithConfigParser sets the configuration parser to be used
func (cfg *Config) WithConfigParser(p utils.Parser) *Config {
	cfg.configParser = p
	return cfg
}

// GetConfigParser returns the configuration parser to be used
func (cfg *Config) GetConfigParser() utils.Parser {
	if cfg.configParser == nil {
		return utils.DefaultConfigParser
	}
	return cfg.configParser
}

// WithDBBuilder sets the database builder to be used
func (cfg *Config) WithDBBuilder(b database.Builder) *Config {
	cfg.buildDb = b
	return cfg
}

// GetDBBuilder returns the database builder to be used
func (cfg *Config) GetDBBuilder() database.Builder {
	if cfg.buildDb == nil {
		return builder.Builder
	}
	return cfg.buildDb
}
