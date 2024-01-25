package partials

import (
	o "operationalcore/components"

	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

type UserFormConfirmPasswordInputProps struct {
	ValidationError string
	Value           string
}

func UserFormConfirmPasswordInput(p *UserFormConfirmPasswordInputProps) g.Node {
	inputProps := &o.InputProps{
		Label:       "Confirm password",
		Name:        "ConfirmPassword",
		InputType:   "password",
		Placeholder: "Confirm password",
		InputProps: []g.Node{
			hx.Post("/users/validate/confirm-password"),
			h.Value(p.Value),
		},
	}

	if p.ValidationError != "" {
		inputProps.HelperText = p.ValidationError
		inputProps.HelperType = o.InputHelperTypeError
	}

	return o.Input(inputProps,
		hx.Target("this"),
		hx.Swap("outerHTML"),
	)
}
