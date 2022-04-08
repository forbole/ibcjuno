package config

import (
	"github.com/forbole/ibcjuno/worker"

	database "github.com/forbole/ibcjuno/db"
	"github.com/forbole/ibcjuno/utils"
)

// GetStartContext sets up database context
func GetStartContext(cfg utils.Config, parseConfig *StartConfig) (*worker.WorkerContext, error) {

	// Create new database context
	databaseCtx := database.NewDatabaseContext(cfg.Database)
	db, err := parseConfig.GetDatabaseBuilder()(databaseCtx)
	if err != nil {
		return nil, err
	}

	return worker.NewWorkerContext(db), nil
}
