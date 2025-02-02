package layout

import (
	"app/internal/components"
	"app/internal/models"
	"app/pkg/reqcontext"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

// icons array
var TopLevelMenuItems = []components.GridMenuItem{
	{
		Icon: "account-group",
		Name: "Users",
		Link: "/users",
		Show: func(permissions models.UserPermissions) bool {
			return permissions.UserAdmin.Access
		},
	},
}

type moduleMenuProps struct {
	Ctx reqcontext.ReqContext
}

func moduleMenu(p *moduleMenuProps) g.Node {

	if p.Ctx.User.UserID == 0 {
		return nil
	}

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
		components.InlineScript("/internal/layout/moduleMenu.js"),
	)
}
