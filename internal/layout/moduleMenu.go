package layout

import (
	"app/internal/components"
	"app/internal/model"
	"app/pkg/reqcontext"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

var AppMenu = []components.GridMenuGroup{
	{
		GroupName: "Admin",
		Items: []components.GridMenuItem{{
			Icon: "account-multiple",
			Name: "Users",
			Link: "/users",
			Show: func(permissions model.UserPermissions) bool {
				return permissions.UserAdmin.Access
			},
		}, {
			Icon: "account-group",
			Name: "Teams",
			Link: "/teams",
			Show: func(permissions model.UserPermissions) bool {
				return permissions.UserAdmin.Access
			},
		}},
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
				Groups:          AppMenu,
				UserPermissions: p.Ctx.User.Permissions,
			}),
		),
		components.InlineScript("/internal/layout/moduleMenu.js"),
	)
}
