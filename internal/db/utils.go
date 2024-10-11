package db

import "database/sql"

// Helper function to handle nullable strings
func NullStringToParam(ns sql.NullString) interface{} {
	if ns.Valid {
		return ns.String
	}
	return nil
}
