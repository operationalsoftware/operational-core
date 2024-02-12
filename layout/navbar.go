package layout

import (
	o "app/components"
	"app/utils"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type navbarProps struct {
	Ctx utils.Context
}

func navbar(p *navbarProps) g.Node {
	return Nav(ID("navbar"),
		o.InlineStyle("/layout/navbar.css"),
		Div(
			Class("nav_links-container"),
			Div(
				Class("logo-container"),
				A(Href("/"),
					Img(
						Alt("Logo"),
						Src("/static/img/logo.png"),
					),
				),
			),
			Div(
				Class("nav_links"),
				g.If(p.Ctx.User.UserID != 0, avatarDropdown(&avatarDropdownProps{
					Ctx: p.Ctx,
				})),
				g.If(p.Ctx.User.UserID != 0, appGallery(&appGalleryProps{
					Ctx: p.Ctx,
				})),
			),
		),
	)
}
