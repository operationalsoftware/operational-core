package handlers

import (
	"net/http"
	"operationalcore/partials"

	"github.com/go-playground/validator/v10"
)

func CreateUserFormLastName(w http.ResponseWriter, r *http.Request) {
	lastName := r.FormValue("last_name")

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Var(lastName, `required,min=3,max=20`)

	var helperText string

	if err != nil {
		helperText = "Last name must be between 3 and 20 characters"
	}

	_ = partials.CreateUserLastNameInput(&partials.CreateUserLastNameInputProps{
		ValidationError: helperText,
		Value:           lastName,
	}).Render(w)

}
