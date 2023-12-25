package handlers

import (
	"net/http"

	"operationalcore/views"

	"github.com/gorilla/mux"
)

func ViewSingleUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_ = views.User(id).Render(w)
}
