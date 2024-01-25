package handlers

import (
	"net/http"
	"operationalcore/partials"

	"github.com/go-playground/validator/v10"
)

func UserFormConfirmPassword(w http.ResponseWriter, r *http.Request) {
	confirmPassword := r.FormValue("ConfirmPassword")
	password := r.FormValue("Password")

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Var(confirmPassword, "required,min=3,max=20")

	var helperText string

	if err != nil {
		helperText = "Passwords must be between 3 and 20 characters"
	}

	if confirmPassword != password {
		helperText = "Passwords do not match"
	}

	_ = partials.UserFormConfirmPasswordInput(&partials.UserFormConfirmPasswordInputProps{
		ValidationError: helperText,
		Value:           confirmPassword,
	}).Render(w)
}
