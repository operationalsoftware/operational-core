package handlers

import (
	"net/http"
	"operationalcore/partials"

	"github.com/go-playground/validator/v10"
)

func CreateUserFormUsername(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Var(username, "required,gte=3,lte=20")

	var helperText string

	if err != nil {
		helperText = "Username must be between 3 and 20 characters"
	}

	_ = partials.CreateUserUsernameInput(&partials.CreateUserUsernameInputProps{
		ValidationError: helperText,
		Value:           username,
	}).Render(w)
}
