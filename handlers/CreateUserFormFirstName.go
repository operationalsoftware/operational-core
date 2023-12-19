package handlers

import (
	"net/http"
	"operationalcore/partials"

	"github.com/go-playground/validator/v10"
)

func CreateUserFormFirstName(w http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("first_name")

	var validate *validator.Validate

	validate = validator.New(validator.WithRequiredStructEnabled())

	err := validate.Var(firstName, "required,gte=3,lte=20")

	var helperText string

	if err != nil {
		helperText = err.Error()
	}

	_ = partials.CreateUserFirstNameInput(&partials.CreateUserFirstNameInputProps{
		ValidationError: helperText,
		Value:           firstName,
	}).Render(w)

}
