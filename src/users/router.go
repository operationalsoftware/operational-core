package users

import (
	"github.com/gorilla/mux"
)

func AddRouter(r *mux.Router) {
	s := r.PathPrefix("/users").Subrouter()

	// Users home
	s.HandleFunc("", indexViewHandler).Methods("GET")

	// Add user
	s.HandleFunc("/add", addUserViewHandler).Methods("GET")
	s.HandleFunc("/add", addUserHandler).Methods("POST")

	// Add API user
	s.HandleFunc("/add-api-user", addAPIUserViewHandler).Methods("GET")
	s.HandleFunc("/add-api-user", addAPIUserHandler).Methods("POST")

	// User page
	s.HandleFunc("/{id}", userViewHandler).Methods("GET")

	// Edit user
	s.HandleFunc("/{id}/edit", editUserViewHandler).Methods("GET")
	s.HandleFunc("/{id}/edit", editUserHandler).Methods("POST")

	// Reset password
	s.HandleFunc("/{id}/reset-password", resetPasswordViewHandler).Methods("GET")
	s.HandleFunc("/{id}/reset-password", resetPasswordHandler).Methods("POST")
}
