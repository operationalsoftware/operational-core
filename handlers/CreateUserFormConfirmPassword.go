package handlers

import (
	"net/http"
	"operationalcore/partials"

	"github.com/go-playground/validator/v10"
)

func CreateUserConfirmPassword(w http.ResponseWriter, r *http.Request) {
	confirmPassword := r.FormValue("confirmPassword")
	password := r.FormValue("password")

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Var(confirmPassword, "required,min=3,max=20")

	var helperText string

	if err != nil {
		helperText = "Passwords must be between 3 and 20 characters"
	}

	if confirmPassword != password {
		helperText = "Passwords do not match"
	}

	_ = partials.CreateUserConfirmPasswordInput(&partials.CreateUserConfirmPasswordInputProps{
		ValidationError: helperText,
		Value:           confirmPassword,
	}).Render(w)

}
