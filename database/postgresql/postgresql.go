package postgresql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	//
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"

	database "github.com/forbole/ibcjuno/database"
	dbtypes "github.com/forbole/ibcjuno/database/types"
	dbutils "github.com/forbole/ibcjuno/database/utils"
	types "github.com/forbole/ibcjuno/types"
	"github.com/forbole/ibcjuno/types/env"
	utils "github.com/forbole/ibcjuno/utils"
)

// Database defines a wrapper around a SQL database and implements functionality
// for data aggregation and exporting.
type Database struct {
	SQL *sqlx.DB
}

// ConnectDatabase creates database connection with configuration set inside
// config.yaml file. It returns a database connection handle or an error
// if the connection fails.
func ConnectDatabase(ctx *database.DatabaseContext) (database.Database, error) {
	dbURI := dbutils.GetEnvOr(env.DatabaseURI, ctx.Cfg.URL)
	dbEnableSSL := dbutils.GetEnvOr(env.DatabaseSSLModeEnable, ctx.Cfg.SSLModeEnable)

	// Configure SSL certificates (optional)
	if dbEnableSSL == "true" {
		dbRootCert := dbutils.GetEnvOr(env.DatabaseSSLRootCert, ctx.Cfg.SSLRootCert)
		dbCert := dbutils.GetEnvOr(env.DatabaseSSLCert, ctx.Cfg.SSLCert)
		dbKey := dbutils.GetEnvOr(env.DatabaseSSLKey, ctx.Cfg.SSLKey)
		dbURI += fmt.Sprintf(" sslmode=require sslrootcert=%s sslcert=%s sslkey=%s",
			dbRootCert, dbCert, dbKey)
	}

	postgresDb, err := sqlx.Open("postgres", dbURI)
	if err != nil {
		return nil, err
	}

	// set max open connections
	postgresDb.SetMaxOpenConns(ctx.Cfg.MaxOpenConnections)
	postgresDb.SetMaxIdleConns(ctx.Cfg.MaxIdleConnections)

	return &Database{
		SQL: postgresDb,
	}, nil
}

// type check to ensure interface is properly implemented
var _ database.Database = &Database{}

// GetTokensPriceID returns the slice of prices id for all tokens stored inside database
func (db *Database) GetTokensPriceID() ([]string, error) {
	var tokens []dbtypes.TokenUnitRow
	var units []string

	query := `SELECT * FROM token_unit`
	err := db.SQL.Select(&tokens, query)
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

	query := `INSERT INTO token_price (token_name, price_id, image, price,
		market_cap, market_cap_rank, fully_diluted_valuation, total_volume,
		high_24h, low_24h, circulating_supply, total_supply, max_supply, 
		ath, atl, timestamp) VALUES`
	var param []interface{}

	for i, ticker := range prices {
		vi := i * 16
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),",
			vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7, vi+8, vi+9, vi+10, vi+11, vi+12, vi+13, vi+14, vi+15, vi+16)
		param = append(param, ticker.Name, ticker.CoingeckoID, ticker.Image, ticker.Price,
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

	_, err := db.SQL.Exec(query, param...)
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
	_, err := db.SQL.Exec(tokenStmt, tokenParams...)
	if err != nil {
		return fmt.Errorf("error while saving tokens: %s", err)
	}

	tokenUnitStmt = tokenUnitStmt[:len(tokenUnitStmt)-1] // Remove trailing ","
	tokenUnitStmt += " ON CONFLICT DO NOTHING"
	_, err = db.SQL.Exec(tokenUnitStmt, tokenUnitParams...)
	if err != nil {
		return fmt.Errorf("error while saving tokens unit: %s", err)
	}

	return err
}

// SaveIBCToken allows to save the given IBC token details inside database
func (db *Database) SaveIBCToken(token types.IBCToken) error {

	// Store IBC token details
	tokenIBCStmt := `INSERT INTO token_ibc (token_name, origin_denom, origin_chain_price_id,
		target_denom, target_chain_price_id, trade_url, timestamp) VALUES `
	var tokenIBCParams []interface{}

	// Initialise token index
	indexIBCToken := 0

	// remove duplicated values retuned from coingecko for tickers
	tickers := dbutils.RemoveDuplicatedTickers(token.Tickers)

	for _, ibc := range tickers {
		cj := indexIBCToken * 7

		tokenIBCStmt += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d),", cj+1, cj+2, cj+3, cj+4, cj+5, cj+6, cj+7)
		tokenIBCParams = append(tokenIBCParams, token.Name, ibc.OriginDenom, ibc.OriginChainPriceID, ibc.TargetDenom,
			ibc.TargetChainPriceID, ibc.TradeURL, ibc.Timestamp)

		// Increase token index
		indexIBCToken++
	}

	tokenIBCStmt = tokenIBCStmt[:len(tokenIBCStmt)-1] // Remove trailing ","
	tokenIBCStmt += `
ON CONFLICT ON CONSTRAINT unique_token_ibc DO UPDATE 
	SET trade_url = excluded.trade_url,
	    timestamp = excluded.timestamp`
	_, err := db.SQL.Exec(tokenIBCStmt, tokenIBCParams...)
	if err != nil {
		return fmt.Errorf("error while saving IBC tokens: %s", err)
	}

	return err
}

// Close implements database.Database
func (db *Database) Close() {
	err := db.SQL.Close()
	if err != nil {
		log.Error().Err(err).Msg("error while closing connection: ")
	}
}
