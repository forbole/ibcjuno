package builder

import (
	database "github.com/forbole/ibcjuno/database"
	"github.com/forbole/ibcjuno/database/postgresql"
)

// DatabaseBuilder represents ConnectDatabase implementation that builds database
// instance based on the configuration set inside config.yaml file
func DatabaseBuilder(ctx *database.DatabaseContext) (database.Database, error) {
	return postgresql.ConnectDatabase(ctx)
}
