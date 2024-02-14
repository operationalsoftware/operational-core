package migrate

import (
	"fmt"
	"log"
	"operationalcore/db"
)

func migrate() bool {
	db := db.UseDB()
	// start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Panic(err)
	}
	defer tx.Rollback()

	//
	// START OF MIGRATION
	//
	var rolesColumnExists bool
	query := `
SELECT EXISTS (
	SELECT 1 FROM sqlite_master WHERE type='table' AND name='User' AND sql LIKE '%Roles%'
)`
	err = db.QueryRow(query).Scan(&rolesColumnExists)
	if err != nil {
		log.Panic(err)
	}

	if !rolesColumnExists {
		var migrations = []string{
			"ALTER TABLE User ADD COLUMN Roles JSON DEFAULT '[]' NOT NULL",
			"ALTER TABLE User ADD COLUMN UserData JSON DEFAULT '{}' NOT NULL",
		}

		for _, migration := range migrations {
			fmt.Printf("Migrating: %s... ", migration)
			_, err = db.Exec(migration)
			if err != nil {
				return false
			}
			fmt.Println("done")
		}
	} else {
		fmt.Println("Database already migrated")
		return true
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

	return true
}
