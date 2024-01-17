package layout

import (
	o "operationalcore/components"
	"operationalcore/utils"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type navbarProps struct {
	Ctx utils.Context
}

func navbar(p *navbarProps) g.Node {
	return Nav(ID("navbar"),
		o.InlineStyle(
			Assets, "/navbar.css",
		),
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
				g.If(p.Ctx.User.UserId != 0, o.AvatarDropdown(&o.AvatarDropdownProps{
					Ctx: p.Ctx,
				})),
				g.If(p.Ctx.User.UserId != 0, o.AppGallery()),
			),
		),
	)
}
