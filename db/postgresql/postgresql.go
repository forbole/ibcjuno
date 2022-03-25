package postgresql

import (
	"database/sql"
	"fmt"

	database "github.com/MonikaCat/ibcjuno/db"
	dbtypes "github.com/MonikaCat/ibcjuno/db/types"
	types "github.com/MonikaCat/ibcjuno/token"
	"github.com/MonikaCat/ibcjuno/token/coingecko"
	utils "github.com/MonikaCat/ibcjuno/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // nolint
	"github.com/rs/zerolog/log"
)

// Database defines a wrapper around a SQL database and implements functionality
// for data aggregation and exporting.
type Database struct {
	Sql  *sql.DB
	Sqlx *sqlx.DB
}

// ConnectDatabase creates database connection with configuration set inside
// config.yaml file. It returns a database connection handle or an error
// if the connection fails.
func ConnectDatabase(ctx *database.DatabaseContext) (database.Database, error) {
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

	// set max open connections
	postgresDb.SetMaxOpenConns(ctx.Cfg.MaxOpenConnections)
	postgresDb.SetMaxIdleConns(ctx.Cfg.MaxIdleConnections)

	return &Database{
		Sql:  postgresDb,
		Sqlx: sqlx.NewDb(postgresDb, "postgresql"),
	}, nil
}

// GetTokensPriceID returns the slice of prices id for all tokens stored inside database
func (db *Database) GetTokensPriceID() ([]string, error) {
	var tokens []dbtypes.TokenUnitRow
	var units []string

	query := `SELECT * FROM token_unit`
	err := db.Sqlx.Select(&tokens, query)
	if err != nil {
		return nil, err
	}

	for _, unit := range tokens {
		if unit.PriceID.String != "" {
			units = append(units, unit.PriceID.String)
		}
	}

	return units, nil
}

// SaveToken allows to save the given token details inside database
func (db *Database) SaveToken(token types.Token) error {
	query := `INSERT INTO token (name) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(query, token.Name)
	if err != nil {
		return err
	}

	query = `INSERT INTO token_unit (token_name, denom, ibc_denom, exponent, price_id) VALUES `
	var params []interface{}

	for i, unit := range token.Units {
		ui := i * 5
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", ui+1, ui+2, ui+3, ui+4, ui+5)
		params = append(params, token.Name, unit.Denom, utils.ToNullString(unit.IBCDenom), unit.Exponent,
			utils.ToNullString(unit.PriceID))
	}

	query = query[:len(query)-1] // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("error while saving token: %s", err)
	}

	return nil
}

// GetTokenPrices allows to query the latest tokens prices from coingecko
func (db *Database) GetTokenPrices() ([]types.TokenPrice, error) {
	// get the list of tokens price id
	ids, err := db.GetTokensPriceID()
	if err != nil {
		log.Error().Err(err).Msg("error while getting tokens price id:")
		return []types.TokenPrice{}, err
	}

	if len(ids) == 0 {
		panic("invalid configuration file: no token price id found inside config.yaml file")
	}

	// get the tokens prices
	prices, err := coingecko.GetLatestTokensPrices(ids)
	if err != nil {
		log.Error().Err(err).Msg("error while getting tokens prices: ")
		return nil, err
	}

	return prices, nil
}

// SaveTokensPrices allows to store the latest tokens price inside database
func (db *Database) SaveTokensPrices(prices []types.TokenPrice) error {
	if len(prices) == 0 {
		return nil
	}

	query := `INSERT INTO token_price (unit_name, price, market_cap, timestamp) VALUES`
	var param []interface{}

	for i, ticker := range prices {
		vi := i * 4
		query += fmt.Sprintf("($%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4)
		param = append(param, ticker.UnitName, ticker.Price, ticker.MarketCap, ticker.Timestamp)
	}

	query = query[:len(query)-1] // Remove trailing ","
	query += `
ON CONFLICT (unit_name) DO UPDATE 
	SET price = excluded.price,
	    market_cap = excluded.market_cap,
	    timestamp = excluded.timestamp
WHERE token_price.timestamp <= excluded.timestamp`

	_, err := db.Sql.Exec(query, param...)
	if err != nil {
		log.Error().Err(err).Msg("error while saving tokens prices: ")
		return err
	}

	return nil
}

// Close implements database.Database
func (db *Database) Close() {
	err := db.Sql.Close()
	if err != nil {
		log.Error().Err(err).Msg("error while closing connection: ")
	}
}
