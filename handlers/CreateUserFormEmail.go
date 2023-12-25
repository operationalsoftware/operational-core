package handlers

import (
	"net/http"
	"operationalcore/partials"

	"github.com/go-playground/validator/v10"
)

func CreateUserFormEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	var validate *validator.Validate

	validate = validator.New(validator.WithRequiredStructEnabled())

	err := validate.Var(email, "required,email")

	var helperText string

	if err != nil {
		helperText = "Email must be a valid email address"
	}

	_ = partials.CreateUserEmailInput(&partials.CreateUserEmailInputProps{
		ValidationError: helperText,
		Value:           email,
	}).Render(w)
}
