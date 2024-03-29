package db

import (
	utils "github.com/forbole/ibcjuno/database/config"
	types "github.com/forbole/ibcjuno/types"
)

type Database interface {

	// Get the slice of prices id for all tokens stored in database
	// Returns error if operation fails
	GetTokensPriceID() ([]string, error)

	// Store given IBC tokens details inside database
	// Returns error if operation fails
	SaveIBCToken(token types.IBCToken) error

	// Store given tokens details inside database
	// Returns error if operation fails
	SaveTokens(token []types.ChainRegistryAsset) error

	// Store latest tokens price in database
	// Returns error if operation fails
	SaveTokensPrices(prices []types.TokenPrice) error

	// Close closes the connection to the database
	Close()
}

// DatabaseContext contains the data used to build a Database instance
type DatabaseContext struct {
	Cfg utils.DatabaseConfig
}

// NewDatabaseContext allows to build new DatabaseContext instance
func NewDatabaseContext(cfg utils.DatabaseConfig) *DatabaseContext {
	return &DatabaseContext{
		Cfg: cfg,
	}
}

// DatabaseBuilder represents a method that allows to build database
type DatabaseBuilder func(ctx *DatabaseContext) (Database, error)
