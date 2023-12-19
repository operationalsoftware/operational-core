package handlers

import (
	"log"
	"net/http"
	"operationalcore/components"
	"operationalcore/db"
	"operationalcore/model"

	"github.com/go-playground/validator/v10"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
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

	var validate *validator.Validate

	validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(user)

	if err != nil {
		_ = components.InputHelper(&components.InputHelperProps{
			Label: "Submission Error",
			Type:  components.InputHelperTypeError,
		}).Render(w)
	}

	// Create user in db
	dbInstance := db.UseDB()
	if dbInstance != nil {
		defer dbInstance.Close()
	}

	if err != nil {
		log.Fatal(err)
	}

	queryErr := model.AddUser(dbInstance, model.User{
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})

	if queryErr != nil {
		log.Fatal(queryErr)
	}

	// Redirect to user view
	w.Header().Set("hx-redirect", "/users")

}
