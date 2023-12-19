package handlers

import (
	"net/http"
	"operationalcore/partials"

	"github.com/go-playground/validator/v10"
)

func CreateUserFormLastName(w http.ResponseWriter, r *http.Request) {
	lastName := r.FormValue("last_name")

	var validate *validator.Validate

	validate = validator.New(validator.WithRequiredStructEnabled())

	err := validate.Var(lastName, "required,gte=3,lte=20")

	var helperText string

	if err != nil {
		helperText = "Last name is required"
	}

	_ = partials.CreateUserLastNameInput(&partials.CreateUserLastNameInputProps{
		ValidationError: helperText,
		Value:           lastName,
	}).Render(w)

}
