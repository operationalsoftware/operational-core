package handlers

import (
	"net/http"
	"operationalcore/partials"

	"github.com/go-playground/validator/v10"
)

func UserFormFirstName(w http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("FirstName")

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Var(firstName, "required,gte=3,lte=20")

	var helperText string

	if err != nil {
		helperText = "First name must be between 3 and 20 characters"
	}

	_ = partials.UserFormFirstNameInput(&partials.UserFormFirstNameInputProps{
		ValidationError: helperText,
		Value:           firstName,
	}).Render(w)
}
