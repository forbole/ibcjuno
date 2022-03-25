package db

import (
	utils "github.com/MonikaCat/ibcjuno/db/config"
	types "github.com/MonikaCat/ibcjuno/token"
)

type Database interface {

	// Get the slice of prices id for all tokens stored in database
	// Returns error if operation fails
	GetTokensPriceID() ([]string, error)

	// Store latest tokens price in database
	// Returns error if operation fails
	SaveTokensPrices(prices []types.TokenPrice) error

	// Query the latest tokens prices from coingecko
	// Returns error if operation fails
	GetTokenPrices() ([]types.TokenPrice, error)

	// Store given token details inside database
	// Returns error if operation fails
	SaveToken(token types.Token) error

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
