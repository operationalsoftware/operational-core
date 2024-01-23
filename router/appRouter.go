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

	// static assets file server
	staticFS := http.FileServer(static.Assets)
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", staticFS))

	// security middleware
	r.Use(middlewares.Security)

	// Authentication middleware
	r.Use(middlewares.Authentication)
	r.Use(middlewares.AuthRedirect)

	// home page
	r.HandleFunc("/", handlers.HomePage).Methods("GET")

	// add subrouters
	AddLoginRouter(r)
	AddUserRouter(r)
	AddLogoutRouter(r)

	// 404
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundPage)

	return r
}
