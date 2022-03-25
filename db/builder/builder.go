package builder

import (
	database "github.com/MonikaCat/ibcjuno/db"
	"github.com/MonikaCat/ibcjuno/db/postgresql"
)

// Builder represents a generic Builder implementation that builds database
// instance based on the configuration set inside config.yaml file
func Builder(ctx *database.Context) (database.Database, error) {
	return postgresql.OpenDB(ctx)
}
