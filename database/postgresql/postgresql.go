package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // nolint
	"github.com/rs/zerolog/log"

	database "github.com/forbole/ibcjuno/database"
	dbtypes "github.com/forbole/ibcjuno/database/types"
	types "github.com/forbole/ibcjuno/types"
	utils "github.com/forbole/ibcjuno/utils"
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

	priceIDList := utils.RemoveDuplicates(units)

	return priceIDList, nil
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
ON CONFLICT DO NOTHING`

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

// SaveIBCTokens allows to save the given IBC tokens details inside database
func (db *Database) SaveIBCTokens(token []types.IBCToken) error {
	// Store tokens name
	tokenStmt := `INSERT INTO token (name) VALUES `
	var tokenParams []interface{}

	// Store token unit details
	tokenUnitStmt := `INSERT INTO token_unit (token_name, denom, exponent, price_id) VALUES `
	var tokenUnitParams []interface{}

	// Store IBC token details
	tokenIBCStmt := `INSERT INTO token_ibc_denom_new (denom, origin_chain, target_denom, target_chain,
		is_stale, trade_url, timestamp) VALUES `
	var tokenIBCParams []interface{}

	// Initialise the indexes
	indexIBCToken := 0
	indexTokenUnits := 0

	for i, ibcDenom := range token {
		u := i * 1
		tokenStmt += fmt.Sprintf("($%d),", u+1)
		tokenParams = append(tokenParams, ibcDenom.Name)

		for _, denomUnit := range ibcDenom.DenomUnits {
			di := indexTokenUnits * 4

			tokenUnitStmt += fmt.Sprintf("($%d,$%d,$%d,$%d),", di+1, di+2, di+3, di+4)

			if denomUnit.Exponent == 0 {
				tokenUnitParams = append(tokenUnitParams, ibcDenom.Name, denomUnit.Denom, denomUnit.Exponent, "")

			} else {
				tokenUnitParams = append(tokenUnitParams, ibcDenom.Name, denomUnit.Denom, denomUnit.Exponent, ibcDenom.CoingeckoID)
			}

			indexTokenUnits++
		}

		for _, ibc := range ibcDenom.Tickers {
			cj := indexIBCToken * 7

			tokenIBCStmt += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d),", cj+1, cj+2, cj+3, cj+4, cj+5, cj+6, cj+7)
			tokenIBCParams = append(tokenIBCParams, ibc.Denom, ibc.OriginChain, ibc.TargetDenom, ibc.TargetChain,
				ibc.IsStale, ibc.TradeURL, ibc.Timestamp)

			indexIBCToken++
		}

	}

	tokenStmt = tokenStmt[:len(tokenStmt)-1] // Remove trailing ","
	tokenStmt += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(tokenStmt, tokenParams...)
	if err != nil {
		return fmt.Errorf("error while saving tokens: %s", err)
	}

	tokenUnitStmt = tokenUnitStmt[:len(tokenUnitStmt)-1] // Remove trailing ","
	tokenUnitStmt += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(tokenUnitStmt, tokenUnitParams...)
	if err != nil {
		return fmt.Errorf("error while saving tokens unit: %s", err)
	}

	tokenIBCStmt = tokenIBCStmt[:len(tokenIBCStmt)-1] // Remove trailing ","
	tokenIBCStmt += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(tokenIBCStmt, tokenIBCParams...)
	if err != nil {
		return fmt.Errorf("error while saving IBC tokens: %s", err)
	}

	log.Info().Msg("** finished processing and storing IBC tokens info in database! ** ")

	return err
}
