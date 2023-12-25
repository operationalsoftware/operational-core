package handlers

import (
	"net/http"

	"operationalcore/views"

	"github.com/gorilla/mux"
)

func EditUserPage(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_ = views.EditUser(id).Render(w)
}
