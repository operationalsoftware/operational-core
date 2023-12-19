package partials

import (
	o "operationalcore/components"

	g "github.com/maragudk/gomponents"
	ghtmx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

type CreateUserFirstNameInputProps struct {
	ValidationError string
	Value           string
}

func CreateUserFirstNameInput(p *CreateUserFirstNameInputProps) g.Node {
	inputProps := &o.InputProps{
		Label:       "First Name",
		Name:        "first_name",
		Placeholder: "Enter first name",
		InputProps: []g.Node{
			ghtmx.Post("/users/create/first-name"),
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
