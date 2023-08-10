package postgresql

import (
	"database/sql"
	"encoding/json"
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

// SaveIBCTokens allows to save the given IBC tokens details inside database
func (db *Database) SaveIBCTokens(token []types.IBCTokenUnit) error {
	// Store tokens name
	tokenStmt := `INSERT INTO token (name) VALUES `
	var tokenparams []interface{}

	// Store token unit details
	tokenUnitStmt := `INSERT INTO token_unit (token_name, denom, exponent, price_id) VALUES `
	var params []interface{}

	// Store IBC token details
	tokenIBCStmt := `INSERT INTO token_ibc_denom (denom, origin_chain, origin_denom, origin_type, 
		symbol, enable, path, channel, counter_party) VALUES `
	var ibcparams []interface{}

	for i, ibcDenom := range token {
		u := i * 1
		tokenStmt += fmt.Sprintf("($%d),", u+1)
		tokenparams = append(tokenparams, ibcDenom.Symbol)

		ui := i * 4
		tokenUnitStmt += fmt.Sprintf("($%d,$%d,$%d,$%d),", ui+1, ui+2, ui+3, ui+4)

		params = append(params, ibcDenom.Symbol, ibcDenom.Denom, ibcDenom.Decimals, utils.ToNullString(ibcDenom.CoingeckoID))

		uj := i * 9
		tokenIBCStmt += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),", uj+1, uj+2, uj+3, uj+4, uj+5, uj+6, uj+7, uj+8, uj+9)
		counterpartyBz, err := json.Marshal(&ibcDenom.CounterParty)
		if err != nil {
			return err
		}
		ibcparams = append(ibcparams, ibcDenom.Denom, ibcDenom.OriginChain, ibcDenom.OriginDenom, ibcDenom.OriginType,
			ibcDenom.Symbol, ibcDenom.Enable, ibcDenom.Path, ibcDenom.Channel, string(counterpartyBz))
	}

	tokenStmt = tokenStmt[:len(tokenStmt)-1] // Remove trailing ","
	tokenStmt += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(tokenStmt, tokenparams...)
	if err != nil {
		return fmt.Errorf("error while saving tokens: %s", err)
	}

	tokenUnitStmt = tokenUnitStmt[:len(tokenUnitStmt)-1] // Remove trailing ","
	tokenUnitStmt += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(tokenUnitStmt, params...)
	if err != nil {
		return fmt.Errorf("error while saving tokens unit: %s", err)
	}

	tokenIBCStmt = tokenIBCStmt[:len(tokenIBCStmt)-1] // Remove trailing ","
	tokenIBCStmt += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(tokenIBCStmt, ibcparams...)
	if err != nil {
		return fmt.Errorf("error while saving IBC tokens: %s", err)
	}

	return err
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
