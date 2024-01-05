package handlers

import (
	"net/http"

	"operationalcore/views"

	"github.com/gorilla/mux"
)

func ViewSingleUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "No ID provided", http.StatusBadRequest)
		return
	}
	_ = views.User(id).Render(w)
}
