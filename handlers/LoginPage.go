package handlers

import (
	"net/http"
	"operationalcore/utils"
	"operationalcore/views"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	_ = views.Login(&views.LoginProps{
		Ctx: ctx,
	}).Render(w)
}
