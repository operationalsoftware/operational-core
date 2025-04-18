package layout

import (
	"app/internal/components"
	"app/pkg/reqcontext"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type avatarMenuProps struct {
	Ctx reqcontext.ReqContext
}

func avatarMenu(p *avatarMenuProps) g.Node {

	if p.Ctx.User.UserID == 0 {
		return nil
	}

	return h.Div(
		h.Class("flex"),

		h.A(
			h.ID("search-button"),
			h.Class("search-button"),
			h.Href("/search"),

			h.Span(
				h.Class("icon"),
				components.Icon(&components.IconProps{
					Identifier: "magnify",
				}),
			),
			h.Span(
				h.Class("search-text"),
				g.Text("Search"),
			),
			h.Span(
				h.Class("shortcut"),
				g.Text("Ctrl + /"),
			),
		),
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
				p.Ctx.User.Email != nil ||
					(p.Ctx.User.FirstName != nil && p.Ctx.User.LastName != nil),
				h.Section(
					(func() g.Node {
						if p.Ctx.User.FirstName == nil || p.Ctx.User.LastName == nil {
							return nil
						}
						return h.P(
							h.Class("name"),
							g.Text(*p.Ctx.User.FirstName+" "+*p.Ctx.User.LastName),
						)
					})(),
					(func() g.Node {
						if p.Ctx.User.Email == nil {
							return nil
						}
						return h.P(
							h.Class("email"),
							g.Text(*p.Ctx.User.Email),
						)
					})(),
				),
			),
			h.A(
				h.Class("logout-btn"),
				h.Href("/auth/logout"),
				g.Text("Logout"),
			),
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
		components.InlineScript("/internal/layout/avatarMenu.js"),
	)
}
