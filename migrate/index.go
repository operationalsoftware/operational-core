package migrate

import (
	"errors"
	"log"
	"operationalcore/db"
)

func InitialiseOrMigrateDB() error {
	db := db.UseDB()
	// start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Panic(err)
	}
	defer tx.Rollback()

	// Check if the initialisation has already been done
	var isInitialised bool
	query := `
SELECT EXISTS (
	SELECT 1 FROM sqlite_master WHERE type='table' AND name='User'
)`
	err = db.QueryRow(query).Scan(&isInitialised)
	if err != nil {
		log.Panic(err)
	}

	if !isInitialised {
		init := initialise()
		if !init {
			return errors.New("initialisation failed")
		} else {
			return nil
		}
	} else {
		migrate := migrate()
		if !migrate {
			return errors.New("migration failed")
		} else {
			return nil
		}
	}
}
