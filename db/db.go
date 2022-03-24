package db

import (
	utils "github.com/MonikaCat/ibcjuno/db/config"
	types "github.com/MonikaCat/ibcjuno/token"
)

type Database interface {

	// Get the slice of prices id for all tokens stored in database
	// Returns error if operation fails
	GetTokensPriceID() ([]string, error)

	// Store the latest tokens price in database
	// Returns error if operation fails
	SaveTokensPrices(prices []types.TokenPrice) error

	// Query the latest tokens prices from coingecko
	// Returns error if operation fails
	GetTokenPrices() ([]types.TokenPrice, error)

	// SaveToken stores the given token details inside database
	// Returns error if operation fails
	SaveToken(token types.Token) error

	// Close closes the connection to the database
	Close()
}

// Context contains the data used to build a Database instance
type Context struct {
	Cfg utils.Config
}

// NewContext allows to build a new Context instance
func NewContext(cfg utils.Config) *Context {
	return &Context{
		Cfg: cfg,
	}
}

// Builder represents a method that allows to build database
type Builder func(ctx *Context) (Database, error)
