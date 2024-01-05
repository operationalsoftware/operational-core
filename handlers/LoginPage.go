package handlers

import (
	"net/http"
	"operationalcore/views"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	_ = views.Login().Render(w)
}
