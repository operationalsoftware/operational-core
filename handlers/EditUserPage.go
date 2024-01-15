package handlers

import (
	"net/http"

	"operationalcore/utils"
	"operationalcore/views"

	"github.com/gorilla/mux"
)

func EditUserPage(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "No ID provided", http.StatusBadRequest)
		return
	}

	ctx := utils.GetContext(r)
	_ = views.EditUser(&views.EditUserProps{
		Id:  id,
		Ctx: ctx,
	}).Render(w)
}
