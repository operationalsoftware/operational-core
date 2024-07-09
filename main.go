package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"app/db"
	"app/migrate"
	"app/routes"
	"app/utils"
)

func main() {
	retcode := 0
	var err error

	defer func() {
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(retcode)
	}()

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

	// Initialise or migrate database
	err = migrate.InitialiseOrMigrateDB()
	if err != nil {
		retcode = 1
		return
	}

	// Initialise some things for start up
	utils.InitCookieInstance()

	// define server
	server := http.Server{
		Addr:    ":3000",
		Handler: routes.Handler(),
	}

	// Bind to a port and pass our router in
	fmt.Println("http://localhost:3000")
	err = server.ListenAndServe()
	if err != nil {
		log.Println(err)
		retcode = 1
		return
	}
}
