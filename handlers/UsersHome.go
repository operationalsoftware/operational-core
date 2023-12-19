package handlers

import (
	"net/http"
	"operationalcore/views"
)

func UsersHome(w http.ResponseWriter, r *http.Request) {
	_ = views.ViewUser().Render(w)
}
