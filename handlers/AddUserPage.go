package handlers

import (
	"net/http"

	"operationalcore/utils"
	"operationalcore/views"
)

func AddUserPage(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	_ = views.AddUser(&views.AddUserProps{
		Ctx: ctx,
	}).Render(w)
}
