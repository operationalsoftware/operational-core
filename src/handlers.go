package src

import (
	"net/http"
	"app/utils"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	_ = indexView(&indexViewProps{
		Ctx: ctx,
	}).Render(w)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	_ = notFoundView(&notFoundViewProps{
		Ctx: ctx,
	}).Render(w)
}
