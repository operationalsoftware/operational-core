package migrate

import (
	"app/db"
	"errors"
	"log"
)

func CheckRequiresInitialisation() bool {
	db := db.UseDB()

	// Check if the initialisation has already been done
	var isInitialised bool
	query := `
SELECT EXISTS (
	SELECT 1 FROM sqlite_master WHERE type='table' AND name='User'
)`
	err := db.QueryRow(query).Scan(&isInitialised)
	if err != nil {
		log.Panic(err)
	}

	return !isInitialised
}

func InitialiseOrMigrateDB() error {
	requireInitialisation := CheckRequiresInitialisation()

	if requireInitialisation {
		init := initialise()
		if !init {
			return errors.New("initialisation failed")
		}
	} else {
		migrate := migrate()
		if !migrate {
			return errors.New("migration failed")
		}
	}

	return nil
}
