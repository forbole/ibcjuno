package worker

import (
	database "github.com/forbole/ibcjuno/database"
)

// Context represents the context that is shared with worker
type Context struct {
	Database database.Database
}

// NewWorkerContext builds new Context instance
func NewWorkerContext(
	db database.Database,
) *Context {
	return &Context{
		Database: db,
	}
}
