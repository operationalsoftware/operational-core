package main

import (
	"log"
	"net/http"
	"os"

	"operationalcore/db"
	"operationalcore/router"

	"github.com/gorilla/handlers"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Connect db
	db.ConnectDB()
	defer db.UseDB().Close()

	r := router.AppRouter()

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":3000", loggedRouter))
}
