package partials

import (
	o "operationalcore/components"

	g "github.com/maragudk/gomponents"
	ghtmx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

type UserFormLastNameInputProps struct {
	ValidationError string
	Value           string
}

func UserFormLastNameInput(p *UserFormLastNameInputProps) g.Node {
	inputProps := &o.InputProps{
		Label:       "Last Name",
		Name:        "LastName",
		Placeholder: "Enter last name",
		InputProps: []g.Node{
			ghtmx.Post("/users/validate/last-name"),
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
