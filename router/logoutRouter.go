package router

import (
	"operationalcore/handlers"

	"github.com/gorilla/mux"
)

func AddLogoutRouter(r *mux.Router) {
	s := r.PathPrefix("/logout").Subrouter()

	s.HandleFunc("", handlers.Logout).Methods("POST")
}
