package partials

import (
	o "operationalcore/components"

	g "github.com/maragudk/gomponents"
	ghtmx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

type UserFormPasswordInputProps struct {
	ValidationError string
	Value           string
}

func UserFormPasswordInput(p *UserFormPasswordInputProps) g.Node {
	inputProps := &o.InputProps{
		Label:       "Password",
		Name:        "Password",
		InputType:   "password",
		Placeholder: "Enter password",
		InputProps: []g.Node{
			ghtmx.Post("/users/validate/password"),
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
