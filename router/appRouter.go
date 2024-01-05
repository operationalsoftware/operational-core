package router

import (
	"net/http"
	"operationalcore/handlers"
	"operationalcore/middlewares"
	"operationalcore/static"

	"github.com/gorilla/mux"
)

func AppRouter() *mux.Router {
	r := mux.NewRouter()

	// TODO: Logging middleware with error handling

	// TODO: general middleware

	// static assets file server
	staticFS := http.FileServer(static.Assets)
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", staticFS))

	// Unprotected  routes
	AddLoginRouter(r)

	// Authentication middleware
	r.Use(middlewares.Authentication)
	r.Use(middlewares.AuthRedirect)

	// protected module routers
	r.HandleFunc("/", handlers.HomePage).Methods("GET")
	AddUserRouter(r)

	// TODO: 404 page
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundPage)

	return r
}
