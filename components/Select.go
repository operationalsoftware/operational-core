package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type Option struct {
	Label g.Node
	Value string
}

func SelectOptions(o []Option, v []string, ID string) g.Node {
	return h.Ul(
		h.Class("dropdown"),
		h.Role("listbox"),
		h.ID(ID),
		g.Group(g.Map(o, func(o Option) g.Node {
			selected := false
			for _, val := range v {
				if val == o.Value {
					selected = true
					break
				}
			}

			return h.Div(
				c.Classes{
					"option":   true,
					"selected": selected,
				},
				h.Role("option"),
				h.DataAttr("value", o.Value),
				h.Span(h.Class("option-label"), o.Label),
				Icon(&IconProps{
					Identifier: "check",
				}),
			)
		})),
	)
}

type SelectProps struct {
	Options       []Option
	Value         []string
	Name          string
	ID            string
	LabelContents g.Node
	Multiple      bool
	Placeholder   string
	Classes       c.Classes
}

func Select(p *SelectProps) g.Node {
	if p.Classes == nil {
		p.Classes = c.Classes{}
	}

	p.Classes["select"] = true

	if p.Name == "" {
		panic("SelectProps.Name must be provided")
	}

	if p.ID == "" {
		p.ID = p.Name
	}

	if p.Placeholder == "" && !p.Multiple {
		p.Placeholder = "Select an option"
	} else if p.Placeholder == "" && p.Multiple {
		p.Placeholder = "Select options"
	}

	dropdownId := p.ID + "-dropdown"

	return h.Label(
		g.If(p.LabelContents != nil, p.LabelContents),
		h.Div(
			p.Classes,
			h.DataAttr("placeholder", p.Placeholder),
			h.DataAttr("name", p.Name),
			g.If(p.Multiple, h.DataAttr("multiple", "true")),
			h.ID(p.ID),
			h.Button(
				h.Role("combobox"),
				g.If(!p.Multiple, h.Aria("label", "select dropdown button")),
				h.Aria("haspopup", "listbox"),
				h.Aria("expanded", "false"),
				h.Aria("controls", dropdownId),

				// A span containing the placeholder/prompt
				h.Span(
					c.Classes{
						"placeholder": true,
						"hidden":      len(p.Value) > 0,
					},
					g.Text(p.Placeholder),
				),

				// A container for the currently selected values
				h.Span(
					h.Class("selected-values"),
					g.If(
						len(p.Value) > 0 && p.Multiple,
						g.Group(g.Map(p.Value, func(v string) g.Node {
							value := v
							var label g.Node
							for _, o := range p.Options {
								if o.Value == v {
									label = o.Label
									break
								}
							}
							return h.Span(
								h.Class("selected-value-chip"),
								h.DataAttr("value", value),
								label,
								h.Span(
									h.Class("remove"),
									g.Text("\u2715"),
								),
								h.Input(
									h.Type("checkbox"),
									h.Name(p.Name),
									h.Value(value),
									h.Checked(),
								),
							)
						})),
					),
					g.If(
						len(p.Value) > 0 && !p.Multiple,
						g.Group(g.Map(p.Value, func(v string) g.Node {
							value := v
							var label g.Node
							for _, o := range p.Options {
								if o.Value == v {
									label = o.Label
									break
								}
							}
							return h.Span(
								h.Class("selected-value"),
								h.DataAttr("value", value),
								label,
								h.Input(
									h.Type("radio"),
									h.Name(p.Name),
									h.Value(value),
									h.Disabled(),
									h.Checked(),
								),
							)
						})),
					),
				),

				// dropdown arrow
				h.Span(
					c.Classes{
						"dropdown-arrow": true,
						"hidden":         len(p.Value) > 0,
					},
					g.Text("\u25BE"), // small downwards pointing triangle
				),

				// clear icon
				h.Span(
					c.Classes{
						"clear-select": true,
						"hidden":       len(p.Value) == 0,
					},
					g.Text("\u2715"), // small close/clear icon
				),
			),
			SelectOptions(p.Options, p.Value, dropdownId),
		),
	)
}
