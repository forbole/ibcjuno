package types

import (
	database "github.com/MonikaCat/ibcjuno/db"
	builder "github.com/MonikaCat/ibcjuno/db/builder"
	utils "github.com/MonikaCat/ibcjuno/utils"
)

// Config contains all the configuration for "start" command
type Config struct {
	configParser utils.Parser
	buildDb      database.Builder
}

// NewConfig allows to build new Config instance
func NewConfig() *Config {
	return &Config{}
}

// GetConfigParser returns the configuration parser
func (cfg *Config) GetConfigParser() utils.Parser {
	if cfg.configParser == nil {
		return utils.DefaultConfigParser
	}
	return cfg.configParser
}

// GetDBBuilder returns the database builder
func (cfg *Config) GetDBBuilder() database.Builder {
	if cfg.buildDb == nil {
		return builder.Builder
	}
	return cfg.buildDb
}
