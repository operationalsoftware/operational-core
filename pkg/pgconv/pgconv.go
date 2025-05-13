package pgconv

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Converts pgtype.Text to *string
func PGTextToStringPtr(pgText pgtype.Text) *string {
	if pgText.Valid {
		return &pgText.String
	}
	return nil
}

// Converts pgtype.Timestamptz to *time.Time
func PGTimestamptzToTimePtr(pgTimestamptz pgtype.Timestamptz) *time.Time {
	if pgTimestamptz.Valid {
		return &pgTimestamptz.Time
	}
	return nil
}

// Converts pgtype.Int4 to *int
func PGInt4ToIntPtr(pgInt pgtype.Int4) *int {
	if pgInt.Valid {
		val := int(pgInt.Int32)
		return &val
	}

	return nil
}
