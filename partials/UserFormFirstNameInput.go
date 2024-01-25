package partials

import (
	o "operationalcore/components"

	g "github.com/maragudk/gomponents"
	ghtmx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

type UserFormFirstNameInputProps struct {
	ValidationError string
	Value           string
}

func UserFormFirstNameInput(p *UserFormFirstNameInputProps) g.Node {
	inputProps := &o.InputProps{
		Label:       "First Name",
		Name:        "FirstName",
		Placeholder: "Enter first name",
		InputProps: []g.Node{
			ghtmx.Post("/users/validate/first-name"),
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
