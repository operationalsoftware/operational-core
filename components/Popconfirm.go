package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type PopconfirmProps struct {
	Id      string
	Icon    string
	Heading string
	Text    string
	Yes     string
	No      string
}

func Popconfirm(p *PopconfirmProps, children ...g.Node) g.Node {
	classes := c.Classes{
		"popconfirm": true,
	}
	return h.Div(
		h.ID(p.Id),
		classes,
		g.Group(children),
		h.Div(
			h.Class("popconfirm-content hide"),
			h.Div(
				h.Class("popconfirm-info"),
				h.P(
					h.Class("popconfirm-icon"),
					Icon(&IconProps{
						Identifier: p.Icon,
					}),
				),
				h.P(
					h.Class("popconfirm-heading"),
					g.Text(p.Heading),
				),
			),
			h.P(
				h.Class("popconfirm-text"),
				g.Text(p.Text),
			),
			h.Div(
				h.Class("popconfirm-actions"),
				Button(&ButtonProps{
					ButtonType: ButtonDanger,
					Size:       ButtonSm,
					Loading:    false,
					Classes: c.Classes{
						"popconfirm-yes": true,
					},
				},
					g.Text(p.Yes),
				),
				Button(&ButtonProps{
					ButtonType: ButtonPrimary,
					Size:       ButtonSm,
					Loading:    false,
					Classes: c.Classes{
						"popconfirm-no": true,
					},
				},
					g.Text(p.No),
				),
			),
		),
		InlineStyle("/Popconfirm.css"),
		InlineScript("/Popconfirm.js"),
	)
}
