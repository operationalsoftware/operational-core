package handlers

import (
	"net/http"

	"operationalcore/utils"
	"operationalcore/views"
)

func CreateUserPage(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	_ = views.CreateUser(&views.CreateUserProps{
		Ctx: ctx,
	}).Render(w)
}
