package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"operationalcore/router"
)


func main() {

	r := router.AppRouter()

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

    // Bind to a port and pass our router in
    log.Fatal(http.ListenAndServe(":3000", loggedRouter))
}

