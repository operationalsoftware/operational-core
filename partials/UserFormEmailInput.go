package partials

import (
	o "operationalcore/components"

	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

type UserFormEmailInputProps struct {
	ValidationError string
	Value           string
}

func UserFormEmailInput(p *UserFormEmailInputProps) g.Node {
	inputProps := &o.InputProps{
		Label:       "Email",
		Name:        "Email",
		InputType:   "email",
		Placeholder: "Enter email",
		InputProps: []g.Node{
			hx.Post("/users/validate/email"),
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
