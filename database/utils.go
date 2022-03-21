package database

import (
	"database/sql"
	"strings"
)

func ToNullString(value string) sql.NullString {
	value = strings.TrimSpace(value)
	return sql.NullString{
		Valid:  value != "",
		String: value,
	}
}
