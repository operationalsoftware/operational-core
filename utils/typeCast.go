package utils

import "database/sql"

// StringToNullString converts a string value to sql.NullString.
func StringToNullString(value string) sql.NullString {
	if value != "" {
		return sql.NullString{String: value, Valid: true}
	}
	return sql.NullString{Valid: false}
}
