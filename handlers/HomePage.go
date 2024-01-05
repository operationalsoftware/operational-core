package handlers

import (
	"fmt"
	"net/http"
	"operationalcore/model"
	"operationalcore/views"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(model.User)
	fmt.Println("HomePage: ", user)
	_ = views.Index().Render(w)
}
