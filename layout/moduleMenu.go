package layout

import (
	"app/components"
	"app/internal/reqcontext"
	"app/models/usermodel"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

// icons array
var TopLevelMenuItems = []components.GridMenuItem{
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

			components.GridMenu(&components.GridMenuProps{
				Items:           TopLevelMenuItems,
				UserPermissions: p.Ctx.User.Permissions,
			}),
		),
		components.InlineScript("/layout/moduleMenu.js"),
	)
}
