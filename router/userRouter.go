package router

import (
	"net/http"
	"operationalcore/handlers"
	"operationalcore/partials"

	"github.com/gorilla/mux"
)

func AddUserRouter(r *mux.Router) {
	s := r.PathPrefix("/users").Subrouter()

	// Users home
	s.HandleFunc("", handlers.UsersPage).Methods("GET")

	// Add user
	s.HandleFunc("/add", handlers.AddUserPage).Methods("GET")
	s.HandleFunc("/add", handlers.AddUser).Methods("POST")

	// partial table
	s.HandleFunc("/table", func(w http.ResponseWriter, r *http.Request) {
		sort := r.URL.Query().Get("sort")
		if sort == "" {
			// /table/sort?Username=ASC
			// set the header to this url
		}
		_ = partials.UsersTable().Render(w)
	}).Methods("GET")

	// User form partials
	s.HandleFunc("/validate/first-name", handlers.UserFormFirstName).Methods("POST")
	s.HandleFunc("/validate/last-name", handlers.UserFormLastName).Methods("POST")
	s.HandleFunc("/validate/email", handlers.UserFormEmail).Methods("POST")
	s.HandleFunc("/validate/username", handlers.UserFormUsername).Methods("POST")
	s.HandleFunc("/validate/password", handlers.UserFormPassword).Methods("POST")
	s.HandleFunc("/validate/confirm-password", handlers.UserFormConfirmPassword).Methods("POST")

	// User page
	s.HandleFunc("/{id}", handlers.UserPage).Methods("GET")

	// Edit user
	s.HandleFunc("/{id}/edit", handlers.EditUserPage).Methods("GET")
	s.HandleFunc("/{id}/edit", handlers.EditUser).Methods("POST")
}
