package components

import (
	"strings"

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
	DefaultText string
	Classes     c.Classes
}

func SelectOptions(o []Option, v []string, name, class, ID string) g.Node {
	return h.Ul(
		h.Class(class),
		h.Role("listbox"),
		h.ID(ID),
		g.Group(g.Map(o, func(o Option) g.Node {
			// checked := "false"
			// for _, val := range v {
			// 	if val == o.Value {
			// 		checked = "true"
			// 		break
			// 	}
			// }
			return h.Li(
				h.Role("option"),
				h.Label(
					g.Text(o.Label),
					h.Input(
						h.Name(name),
						h.Type("radio"),
						h.Value(o.Value),
						// g.Attr("checked", checked),
					),
				),
			)
		})),
	)
}

func MultiSelectOptions(o []Option, v []string, name, class, ID string) g.Node {
	return h.Ul(
		h.Class(class),
		h.Role("listbox"),
		h.ID(ID),
		g.Group(g.Map(o, func(o Option) g.Node {
			checked := "false"
			for _, val := range v {
				if val == o.Value {
					checked = "true"
				}
			}
			return h.Div(
				h.Class("option"),
				h.Li(
					h.Role("option"),
					h.Label(
						h.For(o.Value),
						g.Text(o.Label),
					),
					h.Input(
						h.Name(name),
						h.Type("checkbox"),
						h.Value(o.Value),
						g.Attr("checked", checked),
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

	if p.DefaultText == "" && !p.Multiple {
		p.DefaultText = "Select an option"
	} else if p.DefaultText == "" && p.Multiple {
		p.DefaultText = "Select options"
	}

	dropdownId := p.ID + "-dropdown"

	return h.Div(
		p.Classes,
		h.DataAttr("default-text", p.DefaultText),
		g.If(!p.Multiple, h.DataAttr("default-value", p.Value[0])),
		g.If(p.Multiple, h.DataAttr("default-value", strings.Join(p.Value, ","))),
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
			g.If(!p.Multiple, h.Span(h.Class("default-value"), g.Text(""))),
			g.If(p.Multiple, h.Div(h.Class("selected-values"))),
			h.Span(h.Class("arrow")),
		),
		g.If(!p.Multiple, SelectOptions(p.Options, p.Value, p.Name, "select-dropdown", dropdownId)),
		g.If(p.Multiple, MultiSelectOptions(p.Options, p.Value, p.Name, "multi-select-dropdown", dropdownId)),
		InlineScript("/components/Select.js"),
	)
}
