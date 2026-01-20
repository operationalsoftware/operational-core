package layout

import (
	"app/pkg/reqcontext"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type navbarProps struct {
	ctx reqcontext.ReqContext
}

func navbar(p *navbarProps) g.Node {
	return h.Nav(h.ID("navbar"),
		h.Div(
			h.ID("navbar-logo"),
			h.A(
				h.Href("/"),
				h.Img(
					h.Alt("Logo"),
					h.Src("/static/img/logo.png"),
				),
			),
		),
		h.Div(
			h.ID("navbar-menus"),
			avatarMenu(&avatarMenuProps{
				Ctx: p.ctx,
			}),
			notificationsMenu(&notificationsMenuProps{
				Ctx: p.ctx,
			}),
			moduleMenu(&moduleMenuProps{
				Ctx: p.ctx,
			}),
		),
	)
}
