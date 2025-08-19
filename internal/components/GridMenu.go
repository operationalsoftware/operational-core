package components

import (
	"app/internal/model"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type GridMenuItem struct {
	Icon string
	Name string
	Link string
	Show func(permissions model.UserPermissions) bool
}

type GridMenuGroup struct {
	GroupName string
	Items     []GridMenuItem
}

type GridMenuProps struct {
	Groups          []GridMenuGroup
	UserPermissions model.UserPermissions
}

func GridMenu(p *GridMenuProps) g.Node {

	// filterebased on permissions
	groups := []GridMenuGroup{}

	for _, g := range p.Groups {
		items := []GridMenuItem{}

		for _, i := range g.Items {
			if i.Show != nil {
				show := i.Show(p.UserPermissions)
				if !show {
					continue
				}
			}

			items = append(items, i)
		}

		if len(items) > 0 {
			groups = append(groups, GridMenuGroup{
				GroupName: g.GroupName,
				Items:     items,
			})
		}
	}

	return h.Div(
		h.Class("grid-menu"),
		g.Group(g.Map(groups, func(gr GridMenuGroup) g.Node {

			return h.Div(
				h.Class("grid-menu-group"),

				Divider(g.Text(gr.GroupName)),

				h.Div(
					h.Class("grid-menu-items"),

					g.Group(g.Map(gr.Items, func(i GridMenuItem) g.Node {
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
				),
			)
		})),
	)
}
