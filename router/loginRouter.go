package router

import (
	"operationalcore/handlers"

	"github.com/gorilla/mux"
)

func AddLoginRouter(r *mux.Router) {
	s := r.PathPrefix("/login").Subrouter()

	s.HandleFunc("/password", handlers.LoginPage).Methods("GET")
	s.HandleFunc("/password", handlers.Login).Methods("POST")

}
