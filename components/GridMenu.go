package components

import (
	"app/models/usermodel"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type GridMenuItem struct {
	Icon string
	Name string
	Link string
	Show func(permissions usermodel.UserPermissions) bool
}

type GridMenuProps struct {
	Items           []GridMenuItem
	UserPermissions usermodel.UserPermissions
}

func GridMenu(p *GridMenuProps) g.Node {
	return h.Div(
		h.Class("grid-menu"),
		g.Group(g.Map(p.Items, func(i GridMenuItem) g.Node {

			if i.Show != nil {
				show := i.Show(p.UserPermissions)
				if !show {
					return g.Text("")
				}
			}

			return h.A(
				h.Class("grid-menu-item"),
				h.Href(i.Link),
				Icon(&IconProps{
					Identifier: i.Icon,
				}),
				h.Div(
					h.Class("grid-menu-item-name"),
					g.Text(i.Name),
				),
			)
		})),
	)
}
