package layout

import (
	"app/components"
	reqContext "app/reqcontext"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type avatarDropdownProps struct {
	Ctx reqContext.ReqContext
}

func avatarDropdown(p *avatarDropdownProps) g.Node {

	classes := c.Classes{
		"avatar-dropdown": true,
	}

	return h.Div(
		classes,
		h.Div(
			h.Class("open-btn"),
			components.Icon(&components.IconProps{
				Identifier: "account",
			}),
		),
		h.Div(
			h.Class("dropdown"),
			h.Div(
				h.Class("dropdown-content"),
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
						h.Class("theme-toggle"),
						components.Icon(&components.IconProps{
							Identifier: "theme-system-default",
							Classes: c.Classes{
								"theme-system-default": true,
							},
						}),
						components.Icon(&components.IconProps{
							Identifier: "theme-dark",
							Classes: c.Classes{
								"theme-dark": true,
							},
						}),
						components.Icon(&components.IconProps{
							Identifier: "theme-light",
							Classes: c.Classes{
								"theme-light": true,
							},
						}),
					),
					h.Button(
						h.Class("fullscreen-toggle"),
						components.Icon(&components.IconProps{
							Identifier: "fullscreen",
							Classes: c.Classes{
								"fullscreen": true,
							},
						}),
						components.Icon(&components.IconProps{
							Identifier: "fullscreen-exit",
							Classes: c.Classes{
								"fullscreen-exit": true,
							},
						}),
					),
				),
			),
		),
		components.InlineScript("/layout/avatarDropdown.js"),
	)
}
