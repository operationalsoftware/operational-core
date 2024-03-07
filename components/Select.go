package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type Option struct {
	Label string
	Value string
}

type SelectProps struct {
	Options     []Option
	Value       []string
	Name        string
	ID          string
	Multiple    bool
	Placeholder string
	Classes     c.Classes
}

func SelectOptions(o []Option, v []string, ID string) g.Node {
	return h.Ul(
		h.Class("select-dropdown"),
		h.Role("listbox"),
		h.ID(ID),
		g.Group(g.Map(o, func(o Option) g.Node {
			checked := false
			for _, val := range v {
				if val == o.Value {
					checked = true
					break
				}
			}
			classes := c.Classes{
				"option":  true,
				"checked": checked,
			}

			return h.Div(
				classes,
				h.Role("option"),
				h.DataAttr("value", o.Value),
				h.Span(h.Class("option"), g.Text(o.Label)),
				Icon(&IconProps{
					Identifier: "check",
				}),
			)
		})),
	)
}

func Select(p *SelectProps) g.Node {
	if p.Classes == nil {
		p.Classes = c.Classes{}
	}

	p.Classes["custom-select"] = true

	if p.Name == "" && !p.Multiple {
		p.Name = "select"
	} else if p.Name == "" && p.Multiple {
		p.Name = "multi-select"
	}

	if p.ID == "" && !p.Multiple {
		p.ID = "select"
	} else if p.ID == "" && p.Multiple {
		p.ID = "multi-select"
	}

	if p.Placeholder == "" && !p.Multiple {
		p.Placeholder = "Select an option"
	} else if p.Placeholder == "" && p.Multiple {
		p.Placeholder = "Select options"
	}

	dropdownId := p.ID + "-dropdown"

	placeholderClasses := c.Classes{
		"placeholder": true,
	}

	if len(p.Value) > 0 {
		placeholderClasses["hidden"] = true
	}

	optionClasses := c.Classes{
		"select-value-container": true,
	}

	// get length of value
	if len(p.Value) < 1 {
		optionClasses["hidden"] = true
	}

	return h.Div(
		p.Classes,
		h.DataAttr("placeholder", p.Placeholder),
		h.DataAttr("name", p.Name),
		g.If(p.Multiple, h.DataAttr("multiple", "true")),
		h.ID(p.ID),
		h.Button(
			h.Class("select-button"),
			h.Role("combobox"),
			g.If(!p.Multiple, h.Aria("labelledby", "select button")),
			g.If(p.Multiple, h.Aria("labelledby", "multi-select button")),
			h.Aria("haspopup", "listbox"),
			h.Aria("expanded", "false"),
			h.Aria("controls", dropdownId),
			g.If(p.Multiple, h.Span(
				optionClasses,
				g.Group(g.Map(p.Value, func(v string) g.Node {
					value := v
					label := ""
					for _, o := range p.Options {
						if o.Value == v {
							label = o.Label
							break
						}
					}
					return h.Span(
						h.Class("selected-value"),
						h.DataAttr("value", value),
						g.Text(label),
						h.Span(
							h.Class("remove-icon"),
							g.Text("✕"),
						),
						h.Input(
							h.Type("checkbox"),
							h.Name(p.Name),
							h.Value(value),
							g.Attr("checked", "true"),
						),
					)
				})),
			)),
			h.Span(
				placeholderClasses,
				h.Class("placeholder"),
				g.Text(p.Placeholder),
			),
			h.Span(h.Class("arrow")),
		),
		g.If(!p.Multiple, SelectOptions(p.Options, p.Value, "select-dropdown")),
		g.If(p.Multiple, SelectOptions(p.Options, p.Value, dropdownId)),
		InlineScript("/components/Select.js"),
	)
}

// if p.Multiple {
// 	return h.Span(
// 		h.Class("selected-value"),
// 		h.DataAttr("value", v),
// 		g.Text(v),
// 		h.Span(
// 			h.Class("remove-icon"),
// 			g.Text("✕"),
// 		),
// 		h.Input(
// 			h.Type("checkbox"),
// 			h.Name(p.Name),
// 			h.Value(v),
// 			g.Attr("checked", "true"),
// 		),
// 	)
// }

/*
<span class="selected-value" data-value="1">
Option 1
<span class="remove-icon">✕</span>
<input type="checkbox" name="select" value="1">
</span>
*/
