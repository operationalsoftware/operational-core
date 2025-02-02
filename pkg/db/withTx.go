package db

import (
	"database/sql"
)

func WithTx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	// Start the transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Defer rollback or commit
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()

	// Execute the function within the transaction
	err = fn(tx)
	return err
}
