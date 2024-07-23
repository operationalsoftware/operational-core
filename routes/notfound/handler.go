package notfound

import (
	"app/internal/reqcontext"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	_ = notFoundPage(&notFoundPageProps{
		Ctx: ctx,
	}).Render(w)
	return
}
