package partials

import (
	o "operationalcore/components"

	g "github.com/maragudk/gomponents"
	ghtmx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

type CreateUserConfirmPasswordInputProps struct {
	ValidationError string
	Value           string
}

func CreateUserConfirmPasswordInput(p *CreateUserConfirmPasswordInputProps) g.Node {
	inputProps := &o.InputProps{
		Label:       "Confirm password",
		Name:        "confirmPassword",
		InputType:   "password",
		Placeholder: "Confirm password",
		InputProps: []g.Node{
			ghtmx.Post("/users/create/confirm-password"),
			h.Value(p.Value),
		},
	}

	if p.ValidationError != "" {
		inputProps.HelperText = p.ValidationError
		inputProps.HelperType = o.InputHelperTypeError
	}

	return o.Input(inputProps,
		ghtmx.Target("this"),
		ghtmx.Swap("outerHTML"),
	)
}
