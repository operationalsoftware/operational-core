package partials

import (
	o "operationalcore/components"

	g "github.com/maragudk/gomponents"
	ghtmx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

type UserFormUsernameInputProps struct {
	ValidationError string
	Value           string
}

func UserFormUsernameInput(p *UserFormUsernameInputProps) g.Node {
	inputProps := &o.InputProps{
		Label:       "Username",
		Name:        "Username",
		Placeholder: "Enter username",
		InputProps: []g.Node{
			ghtmx.Post("/users/validate/username"),
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
