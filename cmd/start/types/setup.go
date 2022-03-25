package types

import (
	"github.com/MonikaCat/ibcjuno/worker"

	database "github.com/MonikaCat/ibcjuno/db"
	"github.com/MonikaCat/ibcjuno/utils"
)

// GetStartContext sets up database context
func GetStartContext(cfg utils.Config, parseConfig *StartConfig) (*worker.Context, error) {

	// Create new database context
	databaseCtx := database.NewContext(cfg.Database)
	db, err := parseConfig.GetDBBuilder()(databaseCtx)
	if err != nil {
		return nil, err
	}

	return worker.NewContext(db), nil
}
