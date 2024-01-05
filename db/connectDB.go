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
)

func ConnectDB() {
	once.Do(func() {
		// Connect db
		db, err := sql.Open("sqlite3", "./db/operationalcore.db")
		if err != nil {
			log.Fatal(err)
		}

		err = db.Ping()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Connected to db")

		// Assign the connection to the package-level variable
		dbInstance = db

		// Create table
		statement, _ := dbInstance.Prepare(`
CREATE TABLE IF NOT EXISTS users (
	user_id INTEGER PRIMARY KEY AUTOINCREMENT, 
	username TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL UNIQUE, 
	first_name TEXT NOT NULL, 
	last_name TEXT NOT NULL,
	created DATETIME DEFAULT CURRENT_TIMESTAMP,
	last_login DATETIME NULL,
	hashed_password TEXT NOT NULL,
	failed_login_attempts INTEGER DEFAULT 0,
	login_blocked_until DATETIME NULL 
)`,
		)
		_, err = statement.Exec()
		if err != nil {
			log.Fatal(err)
			statement.Close()
		}
		fmt.Println("Created users table")
		statement.Close()
	})
}

// UseDB returns the shared database connection
func UseDB() *sql.DB {
	return dbInstance
}
