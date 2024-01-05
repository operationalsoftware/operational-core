package handlers

import (
	"fmt"
	"net/http"

	"operationalcore/views"
)

func CreateUserPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateUserPage")
	_ = views.CreateUser().Render(w)
}
