package src

import (
	"net/http"
	"app/assets"
	"app/middlewares"
	"app/src/login"
	"app/src/logout"
	"app/src/users"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()

	// static assets file server
	staticFS := http.FileServer(assets.Assets)
	r.PathPrefix("/static").Handler(staticFS)

	// security middleware
	r.Use(middlewares.Security)

	// Authentication middleware
	r.Use(middlewares.Authentication)
	r.Use(middlewares.AuthRedirect)

	// home page
	r.HandleFunc("/", indexHandler).Methods("GET")

	// add subrouters
	users.AddRouter(r)
	login.AddRouter(r)
	logout.AddRouter(r)

	// 404
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	return r
}
