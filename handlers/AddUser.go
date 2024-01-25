package handlers

import (
	"log"
	"net/http"
	"operationalcore/db"
	"operationalcore/model"
	"operationalcore/utils"
)

func AddUser(w http.ResponseWriter, r *http.Request) {

	// Create user in db
	dbInstance := db.UseDB()

	err := model.AddUser(dbInstance, model.UserToAdd{
		Username:  r.FormValue("Username"),
		Email:     utils.StringToNullString(r.FormValue("Email")),
		FirstName: utils.StringToNullString(r.FormValue("FirstName")),
		LastName:  utils.StringToNullString(r.FormValue("LastName")),
		Password:  r.FormValue("Password"),
	})

	if err != nil {
		log.Panic(err)
	}

	// Redirect to users view
	w.Header().Set("hx-redirect", "/users")
}
