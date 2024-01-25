package handlers

import (
	"fmt"
	"log"
	"net/http"
	"operationalcore/components"
	"operationalcore/db"
	"operationalcore/model"
	"operationalcore/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func EditUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	type User struct {
		FirstName string `validate:"required"`
		LastName  string `validate:"required"`
		Email     string `validate:"email"`
		Username  string `validate:"required,gte=3,lte=20"`
	}

	var user User = User{
		FirstName: r.FormValue("FirstName"),
		LastName:  r.FormValue("LastName"),
		Email:     r.FormValue("Email"),
		Username:  r.FormValue("Username"),
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(user)

	if err != nil {
		_ = components.InputHelper(&components.InputHelperProps{
			Label: "Submission Error",
			Type:  components.InputHelperTypeError,
		}).Render(w)
		return
	}

	dbInstance := db.UseDB()

	if err != nil {
		log.Fatal(err)
	}

	query := model.EditUser(dbInstance, model.User{
		FirstName: utils.StringToNullString(user.FirstName),
		LastName:  utils.StringToNullString(user.LastName),
		Email:     utils.StringToNullString(user.Email),
		Username:  user.Username,
	}, id)

	if query != nil {
		_ = components.InputHelper(&components.InputHelperProps{
			Label: "Submission Error",
			Type:  components.InputHelperTypeError,
		}).Render(w)
		return
	}

	// Redirect to user view
	w.Header().Set("hx-redirect", "/users")

}
