package src

import (
	reqContext "app/reqcontext"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	ctx := reqContext.GetContext(r)
	_ = indexView(&indexViewProps{
		Ctx: ctx,
	}).Render(w)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	ctx := reqContext.GetContext(r)
	_ = notFoundView(&notFoundViewProps{
		Ctx: ctx,
	}).Render(w)
}
