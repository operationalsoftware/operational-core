package partials

import (
	o "operationalcore/components"

	g "github.com/maragudk/gomponents"
	ghtmx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

type CreateUserEmailInputProps struct {
	ValidationError string
	Value           string
}

func CreateUserEmailInput(p *CreateUserEmailInputProps) g.Node {
	inputProps := &o.InputProps{
		Label:       "Email",
		Name:        "email",
		InputType:   "email",
		Placeholder: "Enter email",
		InputProps: []g.Node{
			ghtmx.Post("/users/create/email"),
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
