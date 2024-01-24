package handlers

import (
	"net/http"
	"operationalcore/views"
)

func NotFoundPage(w http.ResponseWriter, r *http.Request) {
	_ = views.Page404().Render(w)
}
