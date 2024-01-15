package handlers

import (
	"net/http"
	"operationalcore/utils"
	"operationalcore/views"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	_ = views.Index(&views.IndexProps{
		Ctx: ctx,
	}).Render(w)
}
