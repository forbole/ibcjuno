package db

import (
	"database/sql"
	"fmt"

	"github.com/MonikaCat/ibc-token/config"
	types "github.com/MonikaCat/ibc-token/token"
	"github.com/MonikaCat/ibc-token/token/coingecko"
	"github.com/MonikaCat/ibc-token/utils"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq" // nolint
	"github.com/rs/zerolog/log"
)

// Database defines a wrapper around a SQL database and implements functionality
// for data aggregation and exporting.
type Database struct {
	*sql.DB
	Sqlx *sqlx.DB
}

// OpenDB opens a database connection with the given database connection info
// from config. It returns a database connection handle or an error if the
// connection fails.
func OpenDB(cfg config.Config) (*Database, error) {
	sslMode := "disable"
	if cfg.DB.SSLMode != "" {
		sslMode = cfg.DB.SSLMode
	}

	schema := "public"
	if cfg.DB.Schema != "" {
		schema = cfg.DB.Schema
	}

	connStr := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s sslmode=%s search_path=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.DB.User, sslMode, schema,
	)

	if cfg.DB.Password != "" {
		connStr += fmt.Sprintf(" password=%s", cfg.DB.Password)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConnections)

	return &Database{DB: db,
		Sqlx: sqlx.NewDb(db, "postgresql"),
	}, nil
}

// GetTokensPriceID returns the slice of price ids for all tokens stored in db
func GetTokensPriceID(db *Database) ([]string, error) {
	var tokens []TokenUnitRow

	query := `SELECT * FROM token_unit`
	err := db.Sqlx.Select(&tokens, query)
	if err != nil {
		return nil, err
	}

	var units []string
	for _, unit := range tokens {
		if unit.PriceID.String != "" {
			units = append(units, unit.PriceID.String)
		}
	}
	fmt.Printf("\n\nunits: %v", units)

	return units, nil
}

// SaveToken allows to save the given token details
func SaveToken(token config.Token, db *Database) error {
	query := `INSERT INTO token (name) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := db.Exec(query, token.Name)
	if err != nil {
		return err
	}

	query = `INSERT INTO token_unit (token_name, denom, exponent, aliases, price_id) VALUES `
	var params []interface{}

	for i, unit := range token.Units {
		ui := i * 5
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", ui+1, ui+2, ui+3, ui+4, ui+5)
		params = append(params, token.Name, unit.Denom, unit.Exponent, pq.StringArray(unit.Aliases),
			utils.ToNullString(unit.PriceID))
	}

	query = query[:len(query)-1] // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"
	_, err = db.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("error while saving token: %s", err)
	}

	return nil
}

// getTokenPrices allows to get the most up-to-date token prices
func GetTokenPrices(db *Database) ([]types.TokenPrice, error) {
	// fmt.Println("LOL")
	// Get the list of tokens price id
	ids, err := GetTokensPriceID(db)
	if err != nil {
		return []types.TokenPrice{}, fmt.Errorf("error while getting tokens price id: %s", err)
	}
	fmt.Printf("Token prices: \n %s", ids)

	if len(ids) == 0 {
		log.Debug().Str("module", "pricefeed").Msg("no traded tokens price id found")
		return []types.TokenPrice{}, nil
	}

	// Get the tokens prices
	prices, err := coingecko.GetTokensPrices(ids)
	if err != nil {
		return nil, fmt.Errorf("error while getting tokens prices: %s", err)
	}

	return prices, nil
}

// SaveTokensPrices allows to save the given prices as the most updated ones
func SaveTokensPrices(prices []types.TokenPrice, db *Database) error {
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

	_, err := db.Exec(query, param...)
	if err != nil {
		return fmt.Errorf("error while saving tokens prices: %s", err)
	}

	return nil
}
