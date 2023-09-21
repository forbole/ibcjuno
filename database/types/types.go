package db

import (
	"database/sql"
)

type TokenUnitRow struct {
	TokenName string         `db:"token_name"`
	Denom     string         `db:"denom"`
	Symbol    string         `db:"symbol"`
	Exponent  int            `db:"exponent"`
	PriceID   sql.NullString `db:"price_id"`
}
