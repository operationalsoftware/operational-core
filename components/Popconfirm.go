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
		classes,
		g.Group(children),
		h.Div(
			h.Class("popconfirm-content"),
			h.ID(p.Id),
			h.Div(
				h.Class("popconfirm-info"),
				h.P(
					h.Class("popconfirm-icon"),
					Icon(p.Icon),
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
					Attributes: []g.Node{
						h.ID("confirmYes"),
					},
				}, g.Text(p.Yes)),
				Button(&ButtonProps{
					ButtonType: ButtonPrimary,
					Size:       ButtonSm,
					Loading:    false,
					Attributes: []g.Node{
						h.ID("confirmNo"),
					},
				}, g.Text(p.No)),
			),
		),
		InlineStyle(Assets, "/Popconfirm.css"),
	)
}
