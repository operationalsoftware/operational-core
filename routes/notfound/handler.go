package notfound

import (
	"app/reqcontext"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	_ = NotFoundPage(&NotFoundPageProps{
		Ctx: ctx,
	}).Render(w)
	return
}
