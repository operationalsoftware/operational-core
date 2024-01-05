package partials

import (
	o "operationalcore/components"

	g "github.com/maragudk/gomponents"
	ghtmx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

type CreateUserPasswordInputProps struct {
	ValidationError string
	Value           string
}

func CreateUserPasswordInput(p *CreateUserPasswordInputProps) g.Node {
	inputProps := &o.InputProps{
		Label:       "Password",
		Name:        "password",
		InputType:   "password",
		Placeholder: "Enter Password",
		InputProps: []g.Node{
			ghtmx.Post("/users/create/password"),
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
