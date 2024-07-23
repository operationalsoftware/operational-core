package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	dbInstance *sql.DB
	once       sync.Once
	err        error
)

func ConnectDB() error {
	once.Do(func() {
		// Connect db
		var db *sql.DB
		db, err = sql.Open("sqlite3", "./app.db")
		if err != nil {
			return
		}

		err = db.Ping()
		if err != nil {
			return
		}

		fmt.Println("Connected to db")

		// Assign the connection to the package-level variable
		dbInstance = db
	})

	if err != nil {
		log.Println("Error connecting to db: ", err)
		return err
	}

	return nil
}

// UseDB returns the shared database connection
func UseDB() *sql.DB {
	return dbInstance
}
