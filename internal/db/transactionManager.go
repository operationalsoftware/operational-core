package db

import (
	"database/sql"
	"log"
)

// TransactionManager manages multiple transactions.
type TransactionManager struct {
	Transactions []*sql.Tx
}

// NewTransactionManager initializes transactions for a variable number of *sql.DB instances.
func NewTransactionManager(dbs ...*sql.DB) (*TransactionManager, error) {
	tm := &TransactionManager{
		Transactions: make([]*sql.Tx, len(dbs)),
	}

	// Begin transactions for each database
	for i, db := range dbs {
		tx, err := db.Begin()
		if err != nil {
			// Rollback any transactions that were started before the error occurred
			for j := 0; j < i; j++ {
				if rbErr := tm.Transactions[j].Rollback(); rbErr != nil {
					log.Printf("Failed to rollback transaction %d: %v", j, rbErr)
				}
			}
			return nil, err
		}
		tm.Transactions[i] = tx
	}

	return tm, nil
}

// Cleanup function commits or rolls back transactions based on error.
func (tm *TransactionManager) Cleanup(err error) {
	if err != nil {
		// Rollback transactions in reverse order
		for i := len(tm.Transactions) - 1; i >= 0; i-- {
			tx := tm.Transactions[i]
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Failed to rollback transaction %d: %v", i, rbErr)
			}
		}
	} else {
		// Commit transactions in reverse order
		for i := len(tm.Transactions) - 1; i >= 0; i-- {
			tx := tm.Transactions[i]
			if commitErr := tx.Commit(); commitErr != nil {
				log.Printf("Failed to commit transaction %d: %v", i, commitErr)
			}
		}
	}
}
