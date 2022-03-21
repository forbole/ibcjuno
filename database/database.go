package database

import (
	"github.com/MonikaCat/ibc-token/database/config"

	types "github.com/MonikaCat/ibc-token/token"
)

type Database interface {
	// Save tokens inside database
	SaveTokens(token types.Token) error

	// Close closes the connection to the database
	Close()
}

// Context contains the config data that is  used to build a Database instance
type Context struct {
	Cfg config.Config
}

// NewContext allows to build a new Context instance
func NewContext(cfg config.Config) *Context {
	return &Context{
		Cfg: cfg,
	}
}

// Builder represents a method that allows to build database
type Builder func(ctx *Context) (Database, error)
