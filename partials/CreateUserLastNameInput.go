package partials

import (
	o "operationalcore/components"

	g "github.com/maragudk/gomponents"
	ghtmx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

type CreateUserLastNameInputProps struct {
	ValidationError string
	Value           string
}

func CreateUserLastNameInput(p *CreateUserLastNameInputProps) g.Node {
	inputProps := &o.InputProps{
		Label:       "Last Name",
		Name:        "last_name",
		Placeholder: "Enter last name",
		InputProps: []g.Node{
			ghtmx.Post("/users/create/last-name"),
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
