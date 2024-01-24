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
			ghtmx.On("click", "htmx.toggleClass(htmx.find('.content-container'), 'show')"),
		),
		h.Div(
			h.Class("content-container"),
			h.Div(
				h.Class("content"),
				h.Div(
					h.Class("user-info"),
					h.P(
						h.Class("name"),
						g.Text(p.Ctx.User.Username),
					),
					g.If(p.Ctx.User.Email.Valid,
						h.P(
							h.Class("email"),
							g.Text(p.Ctx.User.Email.String),
						),
					),
				),
				h.Div(
					h.Class("user-actions"),
					h.P(
						ghtmx.Post("/logout"),
						ghtmx.Swap("none"),
						h.Class("action logout-btn"),
						g.Text("Logout"),
					),
				),
				h.Footer(
					h.Class("footer"),
					h.Div(
						h.Class("footer-content"),
						h.Div(
							h.Class("theme-switcher"),
						),
						h.Div(
							h.Class("fullscreen-switcher"),
						),
					),
				),
			),
		),
		InlineStyle(Assets, "/AvatarDropdown.css"),
		InlineScript(Assets, "/AvatarDropdown.js"),
	)
}
