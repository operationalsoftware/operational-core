package home

import (
	"app/reqcontext"
	"app/routes/home/homeviews"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	_ = homeviews.HomePage(&homeviews.HomePageProps{
		Ctx: reqcontext.GetContext(r),
	}).
		Render(w)

	return
}
