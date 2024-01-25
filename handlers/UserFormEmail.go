package handlers

import (
	"net/http"
	"operationalcore/partials"

	"github.com/go-playground/validator/v10"
)

func UserFormEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("Email")

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Var(email, "email")

	var helperText string

	if err != nil {
		helperText = "Email must be a valid email address"
	}

	_ = partials.UserFormEmailInput(&partials.UserFormEmailInputProps{
		ValidationError: helperText,
		Value:           email,
	}).Render(w)
}
