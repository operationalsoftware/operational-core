package layout

import (
	"app/internal/reqcontext"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type navbarProps struct {
	Ctx reqcontext.ReqContext
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
			g.If(p.Ctx.User.UserID != 0, avatarMenu(&avatarMenuProps{
				Ctx: p.Ctx,
			})),
			g.If(p.Ctx.User.UserID != 0, moduleMenu(&moduleMenuProps{
				Ctx: p.Ctx,
			})),
		),
	)
}
