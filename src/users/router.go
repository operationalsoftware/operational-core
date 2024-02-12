package users

import (
	"net/http"
	"operationalcore/partials"

	"github.com/gorilla/mux"
)

func AddRouter(r *mux.Router) {
	s := r.PathPrefix("/users").Subrouter()

	// Users home
	s.HandleFunc("", indexViewHandler).Methods("GET")

	// Add user
	s.HandleFunc("/add", addUserViewHandler).Methods("GET")
	s.HandleFunc("/add", addUserHandler).Methods("POST")

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
	s.HandleFunc("/validate/first-name", validateFirstNameHandler).Methods("POST")
	s.HandleFunc("/validate/last-name", validateLastNameHandler).Methods("POST")
	s.HandleFunc("/validate/email", validateEmailHandler).Methods("POST")
	s.HandleFunc("/validate/username", validateUsernameHandler).Methods("POST")
	s.HandleFunc("/validate/password", validatePasswordHandler).Methods("POST")
	s.HandleFunc("/validate/confirm-password", validateConfirmPasswordHandler).Methods("POST")

	// User page
	s.HandleFunc("/{id}", userViewHandler).Methods("GET")

	// Edit user
	s.HandleFunc("/{id}/edit", editUserViewHandler).Methods("GET")
	s.HandleFunc("/{id}/edit", editUserHandler).Methods("POST")
}
