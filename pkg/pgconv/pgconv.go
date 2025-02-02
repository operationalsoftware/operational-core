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
