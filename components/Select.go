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
	Options  []Option
	Name     string
	ID       string
	Multiple bool
	Classes  c.Classes
}

func SelectOptions(o []Option, class, ID string) g.Node {
	return h.Ul(
		h.Class(class),
		h.Role("listbox"),
		h.ID(ID),
		g.Group(g.Map(o, func(o Option) g.Node {
			return h.Li(
				h.Role("option"),
				h.Div(
					h.Label(
						h.For(o.Value),
						g.Text(o.Label),
					),
					h.Input(
						h.Type("radio"),
						h.ID(o.Value),
					),
				),
			)
		})),
	)
}

func MultiSelectOptions(o []Option, class, ID string) g.Node {
	return h.Ul(
		h.Class(class),
		h.Role("listbox"),
		h.ID(ID),
		g.Group(g.Map(o, func(o Option) g.Node {
			return h.Div(
				h.Class("option"),
				h.Li(
					h.Role("option"),
					h.Label(
						h.For(o.Value),
						g.Text(o.Label),
					),
					h.Input(
						h.Type("checkbox"),
						h.Value(o.Value),
						g.Attr("checked", "false"),
					),
				),
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

	dropdownId := p.ID + "-dropdown"

	return h.Div(
		p.Classes,
		h.ID(p.ID),
		g.If(p.Multiple, h.DataAttr("multiple", "true")),
		h.Div(
			g.If(!p.Multiple, h.Class("select-button")),
			g.If(p.Multiple, h.Class("multi-select-button")),
			h.Role("combobox"),
			g.If(!p.Multiple, h.Aria("labelledby", "select button")),
			g.If(p.Multiple, h.Aria("labelledby", "multi-select button")),
			h.Aria("haspopup", "listbox"),
			h.Aria("expanded", "false"),
			h.Aria("controls", dropdownId),
			h.Input(
				h.Type("hidden"),
				h.Class("hidden-input"),
				h.Name(p.Name),
				h.Value(""),
			),
			g.If(!p.Multiple, h.Span(h.Class("default-value"), g.Text(""))),
			g.If(p.Multiple, h.Div(h.Class("selected-values"))),
			h.Span(h.Class("arrow")),
		),
		g.If(!p.Multiple, SelectOptions(p.Options, "select-dropdown", dropdownId)),
		g.If(p.Multiple, MultiSelectOptions(p.Options, "multi-select-dropdown", dropdownId)),
		InlineStyle("/components/Select.css"),
		InlineScript("/components/Select.js"),
	)
}
