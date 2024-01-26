package components

import (
	"operationalcore/utils"

	g "github.com/maragudk/gomponents"
	ghtmx "github.com/maragudk/gomponents-htmx"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type AvatarDropdownProps struct {
	Ctx utils.Context
}

func AvatarDropdown(p *AvatarDropdownProps) g.Node {
	classes := c.Classes{
		"avatar-dropdown": true,
	}

	return h.Div(
		classes,
		h.Div(
			h.Class("open-btn"),
			Icon(&IconProps{
				Identifier: "account",
			}),
			// use htmx to toggle class
			ghtmx.On("click", "htmx.toggleClass(htmx.find('.dropdown'), 'show')"),
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
				h.Button(
					h.Class("logout-btn"),
					ghtmx.Post("/logout"),
					ghtmx.Swap("none"),
					g.Text("Logout"),
				),
				h.Section(
					h.Class("actions"),
					h.Button(
						h.Class("theme-switcher"),
					),
					h.Button(
						h.Class("fullscreen-switcher"),
					),
				),
			),
		),
		InlineStyle(Assets, "/AvatarDropdown.css"),
		InlineScript(Assets, "/AvatarDropdown.js"),
	)
}
