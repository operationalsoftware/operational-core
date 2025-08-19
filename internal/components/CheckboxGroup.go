package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type CheckboxOption struct {
	Value string
	Label string
}

type CheckboxGroupProps struct {
	Name    string
	Label   string
	Options []CheckboxOption
	Value   []string
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
		h.Label(h.For(p.Name), g.Text(p.Label)),
		h.Div(
			h.Class("checkbox-options"),
			g.Group(g.Map(p.Options, func(option CheckboxOption) g.Node {
				checked := false
				for _, value := range p.Value {
					if value == option.Value {
						checked = true
						break
					}
				}
				return h.Div(
					h.Class("checkbox-option"),
					Checkbox(&CheckboxProps{
						Name:    p.Name,
						Value:   option.Value,
						Checked: checked,
					}),
					g.Text(option.Label),
				)
			})),
		),
	)
}
