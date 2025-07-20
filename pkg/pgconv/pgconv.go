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

// ----------------------------------------

// Converts *string to pgtype.Text
func StringPtrToPGText(s *string) pgtype.Text {
	if s != nil {
		return pgtype.Text{String: *s, Valid: true}
	}
	return pgtype.Text{Valid: false}
}

// Converts *time.Time to pgtype.Timestamptz
func TimePtrToPGTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t != nil {
		return pgtype.Timestamptz{Time: *t, Valid: true}
	}
	return pgtype.Timestamptz{Valid: false}
}

// Converts *int to pgtype.Int4
func IntPtrToPGInt4(i *int) pgtype.Int4 {
	if i != nil {
		return pgtype.Int4{Int32: int32(*i), Valid: true}
	}
	return pgtype.Int4{Valid: false}
}
