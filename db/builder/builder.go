package builder

import (
	database "github.com/MonikaCat/ibcjuno/db"
	"github.com/MonikaCat/ibcjuno/db/postgresql"
)

// Builder represents a generic Builder implementation that build the proper database
// instance based on the configuration the user has specified
func Builder(ctx *database.Context) (database.Database, error) {
	return postgresql.OpenDB(ctx)
}
