package migrate

import (
	"database/sql"
	"fmt"
)

func checkMigrationRequired(tx *sql.Tx) (bool, error) {

	return false, nil
}

func migrateDB(tx *sql.Tx) error {

	fmt.Println("Migrating DB...")

	var err error

	_, err = tx.Exec(`
ALTER TABLE User ADD COLUMN Permissions JSON DEFAULT '{}' NOT NULL
`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
ALTER TABLE User ADD COLUMN UserData JSON DEFAULT '{}' NOT NULL
`)
	if err != nil {
		return err
	}

	fmt.Println("...done")

	return nil
}
