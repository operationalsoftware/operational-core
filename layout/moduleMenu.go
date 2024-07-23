package layout

import (
	"app/components"
	"app/internal/reqcontext"
	"app/models/usermodel"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type ModuleGridItem struct {
	Icon string
	Name string
	Link string
	Show func(usermodel.UserPermissions) bool
}

// icons array
var ModuleGridItems = []ModuleGridItem{
	{
		Icon: "account-group",
		Name: "Users",
		Link: "/users",
		Show: func(permissions usermodel.UserPermissions) bool {
			return permissions.UserAdmin.Access
		},
	},
}

type moduleMenuProps struct {
	Ctx reqcontext.ReqContext
}

func moduleMenu(p *moduleMenuProps) g.Node {

	return h.Div(
		h.Button(
			h.ID("navbar-module-menu-button"),
			h.Class("menu-button"),
			components.Icon(&components.IconProps{
				Identifier: "dots",
			}),
		),
		// styled with position: fixed
		h.Div(
			h.Class("dropdown-panel"),
			h.ID("navbar-module-menu"),
			g.Group(g.Map(ModuleGridItems, func(i ModuleGridItem) g.Node {
				if i.Show != nil {
					show := i.Show(p.Ctx.User.Permissions)
					if !show {
						return g.Text("")
					}
				}

				return h.A(
					h.Class("item"),
					h.Href(i.Link),
					components.Icon(&components.IconProps{
						Identifier: i.Icon,
					}),
					h.Div(
						h.Class("name"),
						g.Text(i.Name),
					),
				)
			})),
		),
		components.InlineScript("/layout/moduleMenu.js"),
	)
}
