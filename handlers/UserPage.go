package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"operationalcore/utils"
	"operationalcore/views"

	"github.com/gorilla/mux"
)

func UserPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	ctx := utils.GetContext(r)
	_ = views.User(&views.UserProps{
		Id:  id,
		Ctx: ctx,
	}).Render(w)
}
