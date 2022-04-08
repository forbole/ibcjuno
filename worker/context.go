package worker

import (
	database "github.com/forbole/ibcjuno/db"
)

// WorkerContext represents the context that is shared with worker
type WorkerContext struct {
	Database database.Database
}

// NewWorkerContext builds new WorkerContext instance
func NewWorkerContext(
	db database.Database,
) *WorkerContext {
	return &WorkerContext{
		Database: db,
	}
}
