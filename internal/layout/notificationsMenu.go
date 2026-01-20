package layout

import (
	"app/internal/components"
	"app/pkg/reqcontext"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type notificationsMenuProps struct {
	Ctx reqcontext.ReqContext
}

func notificationsMenu(p *notificationsMenuProps) g.Node {
	if p.Ctx.User.UserID == 0 {
		return nil
	}

	return h.Div(
		h.Class("notifications-menu"),
		h.Button(
			h.ID("navbar-notifications-menu-button"),
			h.Class("menu-button"),
			h.Aria("label", "Notifications"),
			components.Icon(&components.IconProps{
				Identifier: "bell-outline",
			}),
			h.Span(h.Class("notifications-badge")),
		),
		h.Div(
			h.Class("dropdown-panel"),
			h.ID("navbar-notifications-menu"),
			h.Div(
				h.Class("notifications-tray-loading"),
				g.Text("Loading notifications..."),
			),
		),
		components.InlineScript("/internal/layout/notificationsMenu.js"),
	)
}
