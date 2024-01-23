package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"operationalcore/db"
	"operationalcore/router"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

// TODO: Graceful shutdown

func loadEnv() {
	// Check if GO_ENV is "staging" or "production"
	goEnv := os.Getenv("GO_ENV")
	if goEnv != "staging" && goEnv != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func main() {
	// Connect db
	db.ConnectDB()
	defer db.UseDB().Close()

	loadEnv()

	r := router.AppRouter()

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	// Bind to a port and pass our router in
	fmt.Println("http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", loggedRouter))
}
