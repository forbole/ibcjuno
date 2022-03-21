package postgresql

import (
	"database/sql"
	"fmt"

	database "github.com/MonikaCat/ibc-token/database"
	types "github.com/MonikaCat/ibc-token/token"
	"github.com/lib/pq"
)

// Builder creates a database connection with the given database connection info
// from config. It returns a database connection handle or an error if the
// connection fails.
func Builder(ctx *database.Context) (database.Database, error) {
	sslMode := "disable"
	if ctx.Cfg.SSLMode != "" {
		sslMode = ctx.Cfg.SSLMode
	}

	schema := "public"
	if ctx.Cfg.Schema != "" {
		schema = ctx.Cfg.Schema
	}

	connStr := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s sslmode=%s search_path=%s",
		ctx.Cfg.Host, ctx.Cfg.Port, ctx.Cfg.Name, ctx.Cfg.User, sslMode, schema,
	)

	if ctx.Cfg.Password != "" {
		connStr += fmt.Sprintf(" password=%s", ctx.Cfg.Password)
	}

	postgresDb, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	fmt.Printf("host=%s port=%d dbname=%s user=%s sslmode=%s search_path=%s",
		ctx.Cfg.Host, ctx.Cfg.Port, ctx.Cfg.Name, ctx.Cfg.User, sslMode, schema)
	// Set max open connections
	postgresDb.SetMaxOpenConns(ctx.Cfg.MaxOpenConnections)
	postgresDb.SetMaxIdleConns(ctx.Cfg.MaxIdleConnections)

	return &Database{
		Sql: postgresDb,
	}, nil
}

// type check to ensure interface is properly implemented
var _ database.Database = &Database{}

// Database defines a wrapper around a SQL database and implements functionality
// for data aggregation and exporting.
type Database struct {
	Sql *sql.DB
}

// SaveTokens allows to save the given token details
func (db *Database) SaveTokens(token types.Token) error {
	query := `INSERT INTO token (name) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(query, token.Name)
	if err != nil {
		return err
	}

	query = `INSERT INTO token_unit (token_name, denom, exponent, aliases, price_id) VALUES `
	var params []interface{}

	for i, unit := range token.Units {
		ui := i * 5
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", ui+1, ui+2, ui+3, ui+4, ui+5)
		params = append(params, token.Name, unit.Denom, unit.Exponent, pq.StringArray(unit.Aliases),
			database.ToNullString(unit.PriceID))
	}

	query = query[:len(query)-1] // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("error while saving token: %s", err)
	}

	return nil
}

// Close implements database.Database
func (db *Database) Close() {
	err := db.Sql.Close()
	if err != nil {
		fmt.Errorf("error while closing connection", "err", err)
	}
}
