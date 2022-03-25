package types

import (
	database "github.com/MonikaCat/ibcjuno/db"
	builder "github.com/MonikaCat/ibcjuno/db/builder"
	utils "github.com/MonikaCat/ibcjuno/utils"
)

// StartConfig contains all the configuration for "start" command
type StartConfig struct {
	configParser utils.ConfigParser
	buildDb      database.DatabaseBuilder
}

// NewStartConfig allows to build new StartConfig instance
func NewStartConfig() *StartConfig {
	return &StartConfig{}
}

// GetConfigParser returns the configuration parser
func (cfg *StartConfig) GetConfigParser() utils.ConfigParser {
	if cfg.configParser == nil {
		return utils.DefaultConfigParser
	}
	return cfg.configParser
}

// GetDBBuilder returns the database builder
func (cfg *StartConfig) GetDatabaseBuilder() database.DatabaseBuilder {
	if cfg.buildDb == nil {
		return builder.DatabaseBuilder
	}
	return cfg.buildDb
}
