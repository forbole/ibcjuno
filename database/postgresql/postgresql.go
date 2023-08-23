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

	query := `INSERT INTO token_price (price_id, token_name, image, price,
		market_cap, market_cap_rank, fully_diluted_valuation, total_volume,
		high_24h, low_24h, circulating_supply, total_supply, max_supply, 
		ath, atl, timestamp) VALUES`
	var param []interface{}

	for i, ticker := range prices {
		vi := i * 16
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),",
			vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7, vi+8, vi+9, vi+10, vi+11, vi+12, vi+13, vi+14, vi+15, vi+16)
		param = append(param, ticker.CoingeckoID, ticker.Name, ticker.Image, ticker.Price,
			ticker.MarketCap, ticker.MarketCapRank, ticker.FullyDilutedValuation, ticker.TotalVolume,
			ticker.High24Hrs, ticker.Low24Hrs, ticker.CirculatingSupply, ticker.TotalSupply,
			ticker.MaxSupply, ticker.ATH, ticker.ATL, ticker.Timestamp)
	}

	query = query[:len(query)-1] // Remove trailing ","
	query += `
ON CONFLICT (price_id) DO UPDATE 
		SET token_name = excluded.token_name,
			image = excluded.image, 
			price = excluded.price, 
			market_cap = excluded.market_cap, 
			market_cap_rank = excluded.market_cap_rank, 
			fully_diluted_valuation = excluded.fully_diluted_valuation, 
			total_volume = excluded.total_volume, 
			high_24h = excluded.high_24h, 
			low_24h = excluded.low_24h, 
			circulating_supply = excluded.circulating_supply, 
			total_supply = excluded.total_supply, 
			max_supply = excluded.max_supply, 
			ath = excluded.ath, 
			atl = excluded.atl, 
			timestamp = excluded.timestamp`

	_, err := db.Sql.Exec(query, param...)
	if err != nil {
		log.Error().Err(err).Msg("error while saving tokens prices: ")
		return err
	}

	log.Info().Msg("*** Finished storing token prices in database ***")

	return nil
}

// SaveTokens allows to save the given tokens details inside database
func (db *Database) SaveTokens(token []types.ChainRegistryAsset) error {
	// Store tokens name
	tokenStmt := `INSERT INTO token (name) VALUES `
	var tokenParams []interface{}

	// Store token unit details
	tokenUnitStmt := `INSERT INTO token_unit (token_name, symbol, denom, exponent, price_id) VALUES `
	var tokenUnitParams []interface{}

	indexTokenUnits := 0

	for i, ibcDenom := range token {
		u := i * 1
		tokenStmt += fmt.Sprintf("($%d),", u+1)
		tokenParams = append(tokenParams, ibcDenom.Name)

		for _, denomUnit := range ibcDenom.DenomUnits {
			di := indexTokenUnits * 5

			tokenUnitStmt += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", di+1, di+2, di+3, di+4, di+5)

			if denomUnit.Exponent > 0 {
				tokenUnitParams = append(tokenUnitParams, ibcDenom.Name, ibcDenom.Symbol, denomUnit.Denom, denomUnit.Exponent, ibcDenom.CoingeckoID)

			} else {
				tokenUnitParams = append(tokenUnitParams, ibcDenom.Name, ibcDenom.Symbol, denomUnit.Denom, denomUnit.Exponent, "")
			}

			indexTokenUnits++
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

	return err
}

// SaveIBCTokens allows to save the given IBC tokens details inside database
func (db *Database) SaveIBCTokens(token []types.IBCToken) error {

	// Store IBC token details
	tokenIBCStmt := `INSERT INTO token_ibc (token_name, origin_denom, origin_chain_price_id,
		target_denom, target_chain_price_id, trade_url, timestamp) VALUES `
	var tokenIBCParams []interface{}

	// Initialise token index
	indexIBCToken := 0

	for _, ibcDenom := range token {
		for _, ibc := range ibcDenom.Tickers {
			cj := indexIBCToken * 7

			tokenIBCStmt += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d),", cj+1, cj+2, cj+3, cj+4, cj+5, cj+6, cj+7)
			tokenIBCParams = append(tokenIBCParams, ibcDenom.Name, ibc.OriginDenom, ibc.OriginChainPriceID, ibc.TargetDenom,
				ibc.TargetChainPriceID, ibc.TradeURL, ibc.Timestamp)

			// Increase token index
			indexIBCToken++
		}
	}

	tokenIBCStmt = tokenIBCStmt[:len(tokenIBCStmt)-1] // Remove trailing ","
	tokenIBCStmt += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(tokenIBCStmt, tokenIBCParams...)
	if err != nil {
		return fmt.Errorf("error while saving IBC tokens: %s", err)
	}

	return err
}

// Close implements database.Database
func (db *Database) Close() {
	err := db.Sql.Close()
	if err != nil {
		log.Error().Err(err).Msg("error while closing connection: ")
	}
}
