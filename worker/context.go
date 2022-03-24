package worker

import (
	database "github.com/MonikaCat/ibcjuno/db"
	
)

// Context represents the context that is shared among different workers
type Context struct {
	Database database.Database
}

// NewContext builds a new Context instance
func NewContext(
	db database.Database,
) *Context {
	return &Context{
		Database: db,
	}
}
