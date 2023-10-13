package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

func addContactsRouter(r *mux.Router) {
	s := r.PathPrefix("/contacts").Subrouter()

	s.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Contacts!\n"))
	}).Methods("GET")
}
