package layout

import (
	"app/components"
	"app/internal/reqcontext"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type avatarMenuProps struct {
	Ctx reqcontext.ReqContext
}

func avatarMenu(p *avatarMenuProps) g.Node {

	return h.Div(
		h.Button(
			h.ID("navbar-avatar-menu-button"),
			h.Class("menu-button"),
			components.Icon(&components.IconProps{
				Identifier: "account",
			}),
		),
		// styled with position: fixed
		h.Div(
			h.Class("dropdown-panel"),
			h.ID("navbar-avatar-menu"),
			h.Section(
				h.Class("username"),
				g.Text(p.Ctx.User.Username),
			),
			g.If(
				p.Ctx.User.Email.Valid || (p.Ctx.User.FirstName.Valid && p.Ctx.User.LastName.Valid),
				h.Section(
					g.If(
						p.Ctx.User.FirstName.Valid && p.Ctx.User.LastName.Valid,
						h.P(
							h.Class("name"),
							g.Text(p.Ctx.User.FirstName.String+" "+p.Ctx.User.LastName.String),
						),
					),
					g.If(
						p.Ctx.User.Email.Valid,
						h.P(
							h.Class("email"),
							g.Text(p.Ctx.User.Email.String),
						),
					),
				),
			),
			h.Button(h.Class("logout-btn"), g.Text("Logout")),
			h.Section(
				h.Class("actions"),
				h.Button(
					h.ID("theme-toggle-button"),
					components.Icon(&components.IconProps{
						Identifier: "theme-system-default",
						Classes: c.Classes{
							"theme-system-default-icon": true,
						},
					}),
					components.Icon(&components.IconProps{
						Identifier: "theme-dark",
						Classes: c.Classes{
							"theme-dark-icon": true,
						},
					}),
					components.Icon(&components.IconProps{
						Identifier: "theme-light",
						Classes: c.Classes{
							"theme-light-icon": true,
						},
					}),
				),
				h.Button(
					h.ID("fullscreen-toggle-button"),
					components.Icon(&components.IconProps{
						Identifier: "fullscreen",
						Classes: c.Classes{
							"fullscreen-icon": true,
						},
					}),
					components.Icon(&components.IconProps{
						Identifier: "fullscreen-exit",
						Classes: c.Classes{
							"fullscreen-exit-icon": true,
						},
					}),
				),
			),
		),
		components.InlineScript("/layout/avatarMenu.js"),
	)
}
