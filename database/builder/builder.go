package builder

import (
	"github.com/MonikaCat/ibc-token/database"

	"github.com/MonikaCat/ibc-token/database/postgresql"
)

// Builder represents a generic Builder implementation that build the proper database
// instance based on the configuration the user has specified
func Builder(ctx *database.Context) (database.Database, error) {
	return postgresql.Builder(ctx)
}
