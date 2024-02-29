package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type CheckboxOption struct {
	Value   string
	Label   string
	Checked bool
}

type CheckboxGroupProps struct {
	Name    string
	Label   string
	Options []CheckboxOption
	Classes c.Classes
}

func CheckboxGroup(p *CheckboxGroupProps) g.Node {
	if p.Options == nil {
		p.Options = []CheckboxOption{}
	}

	if p.Name == "" {
		p.Name = "checkbox-group"
	}

	if p.Label == "" {
		p.Label = "Checkbox Group"
	}

	if p.Classes == nil {
		p.Classes = c.Classes{}
	}

	p.Classes["checkbox-group"] = true

	return h.Div(
		p.Classes,
		h.Input(
			h.Type("hidden"),
			h.Class("hidden-input"),
			h.Name(p.Name),
			h.Value(""),
		),
		InputLabel(&InputLabelProps{
			For: p.Name,
		}, g.Text(p.Label)),
		h.Div(
			h.Class("checkbox-options"),
			g.Group(g.Map(p.Options, func(option CheckboxOption) g.Node {
				return h.Div(
					h.Class("checkbox-option"),
					Checkbox(&CheckboxProps{
						Name:    option.Label,
						Value:   option.Value,
						Checked: option.Checked,
					}),
					g.Text(option.Label),
				)
			})),
		),
		InlineScript("/components/CheckboxGroup.js"),
		InlineStyle("/components/CheckboxGroup.css"),
	)
}
