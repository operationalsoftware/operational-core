package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"operationalcore/db"
	"operationalcore/migrate"
	"operationalcore/src"
	"operationalcore/utils"

	"github.com/gorilla/handlers"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	retcode := 0
	defer func() { os.Exit(retcode) }()

	var err error

	// Load environment (if not in production or staging)
	err = utils.LoadEnv()
	if err != nil {
		retcode = 1
		return
	}

	// Verify environment variables
	err = utils.VerifyEnv()
	if err != nil {
		retcode = 1
		return
	}

	// Create database connection
	err = db.ConnectDB()
	if err != nil {
		retcode = 1
		return
	}
	defer db.UseDB().Close()

	// Run migrations
	err = migrate.Initialise()
	if err != nil {
		retcode = 1
		return
	}

	// Initialise some things for start up
	utils.InitCookieInstance()

	// Get router
	r := src.Router()

	// Wrap router with Gorilla logging
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	// Bind to a port and pass our router in
	fmt.Println("http://localhost:3000")
	err = http.ListenAndServe(":3000", loggedRouter)
	if err != nil {
		log.Println(err)
		retcode = 1
		return
	}
}
