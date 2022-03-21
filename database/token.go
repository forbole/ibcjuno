package database

import (
	"fmt"

	types "github.com/MonikaCat/ibc-token/token"
	"github.com/lib/pq"
)

// SaveTokens allows to save the given token details
func (db *Database) SaveTokens(token types.Token) error {
	query := `INSERT INTO token (name) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := db.Sqlx.Exec(query, token.Name)
	if err != nil {
		return err
	}

	query = `INSERT INTO token_unit (token_name, denom, exponent, aliases, price_id) VALUES `
	var params []interface{}

	for i, unit := range token.Units {
		ui := i * 5
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", ui+1, ui+2, ui+3, ui+4, ui+5)
		params = append(params, token.Name, unit.Denom, unit.Exponent, pq.StringArray(unit.Aliases),
			ToNullString(unit.PriceID))
	}

	query = query[:len(query)-1] // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"
	_, err = db.Sqlx.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("error while saving token: %s", err)
	}

	return nil
}
