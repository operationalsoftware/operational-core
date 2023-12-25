package handlers

import (
	"log"
	"net/http"
	"operationalcore/components"
	"operationalcore/db"
	"operationalcore/model"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func EditUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		log.Fatal("No id provided")
		return
	}

	type User struct {
		FirstName string `validate:"required"`
		LastName  string `validate:"required"`
		Email     string `validate:"required,email"`
		Username  string `validate:"required,gte=3,lte=20"`
	}

	var user User = User{
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
		Email:     r.FormValue("email"),
		Username:  r.FormValue("username"),
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(user)

	if err != nil {
		_ = components.InputHelper(&components.InputHelperProps{
			Label: "Submission Error",
			Type:  components.InputHelperTypeError,
		}).Render(w)
		return
	}

	dbInstance := db.UseDB()
	// defer dbInstance.Close()

	if err != nil {
		log.Fatal(err)
	}

	query := model.EditUser(dbInstance, model.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
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
