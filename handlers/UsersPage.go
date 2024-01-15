package handlers

import (
	"net/http"
	"operationalcore/utils"
	"operationalcore/views"
)

func UsersPage(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	_ = views.Users(&views.UsersProps{
		Ctx: ctx,
	}).Render(w)
}
