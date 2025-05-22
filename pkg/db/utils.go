package db

import (
	"database/sql"
	"time"
)

// Helper function to handle nullable strings
func NullStringToParam(ns sql.NullString) interface{} {
	if ns.Valid {
		return ns.String
	}
	return (*string)(nil)
}

// Helper function to handle nullable time
func NullTimeToParam(nt sql.NullTime) interface{} {
	if nt.Valid {
		return nt.Time
	}
	return (*time.Time)(nil)
}
