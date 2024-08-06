package migrate

import (
	"app/internal/db"
)

func MigrateDB() error {

	conn := db.UseDB()
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	initialisationRequired, err := checkInitialisationRequired(tx)
	if err != nil {
		return err
	}

	if initialisationRequired {
		err := initialiseDB(tx)
		if err != nil {
			return err
		}
	}

	migrationRequired, err := checkMigrationRequired(tx)
	if migrationRequired {
		err := migrateDB(tx)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
