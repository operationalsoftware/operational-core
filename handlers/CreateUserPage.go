package handlers

import (
	"net/http"

	"operationalcore/views"
)

func CreateUserPage(w http.ResponseWriter, r *http.Request) {
	_ = views.CreateUser().Render(w)
}
