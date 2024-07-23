package home

import (
	"app/internal/reqcontext"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	ctx := reqcontext.GetContext(r)
	_ = homePage(&homePageProps{
		Ctx: ctx,
	}).Render(w)
	return

}
