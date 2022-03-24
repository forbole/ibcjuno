package db

import (
	"database/sql"
	"time"
)

type TokenUnitRow struct {
	TokenName string         `db:"token_name"`
	Denom     string         `db:"denom"`
	IBCDenom  sql.NullString `db:"ibc_denom"`
	Exponent  int            `db:"exponent"`
	PriceID   sql.NullString `db:"price_id"`
}

type TokenRow struct {
	Name       string `db:"name"`
	TradedUnit string `db:"traded_unit"`
}

// --------------------------------------------------------------------------------------------------------------------

// TokenPriceRow represent a row of the table token_price inside database
type TokenPriceRow struct {
	ID        string    `db:"id"`
	Name      string    `db:"unit_name"`
	Price     float64   `db:"price"`
	MarketCap int64     `db:"market_cap"`
	Timestamp time.Time `db:"timestamp"`
}

// NewTokenPriceRow allows to create NewTokenPriceRow
func NewTokenPriceRow(name string, currentPrice float64, marketCap int64, timestamp time.Time) TokenPriceRow {
	return TokenPriceRow{
		Name:      name,
		Price:     currentPrice,
		MarketCap: marketCap,
		Timestamp: timestamp,
	}
}

// Returns true if u and v represent the same row
func (u TokenPriceRow) Equals(v TokenPriceRow) bool {
	return u.Name == v.Name &&
		u.Price == v.Price &&
		u.MarketCap == v.MarketCap &&
		u.Timestamp.Equal(v.Timestamp)
}
