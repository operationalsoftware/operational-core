package layout

import (
	"app/components"
	"app/utils"

	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type Item struct {
	Icon string
	Name string
	Link string
	Role []string
}

type appGalleryProps struct {
	Ctx utils.Context
}

func appGallery(p *appGalleryProps) g.Node {
	classes := c.Classes{
		"app-gallery": true,
	}

	// icons array
	icons := []Item{
		{
			Icon: "account-group",
			Name: "Users",
			Link: "/users",
			Role: []string{"User Admin"},
		},
		{
			Icon: "github",
			Name: "Github",
			Link: "/github",
			Role: []string{"User", "Dummy"},
		},
	}

	return h.Div(
		classes,
		h.Div(
			h.Class("app-gallery__button"),
			components.Icon(&components.IconProps{
				Identifier: "dots",
			}),
			// use htmx to toggle class
			hx.On("click", "htmx.toggleClass(htmx.find('.app-gallery-content__container'), 'show')"),
		),
		h.Div(
			h.Class("app-gallery-content__container"),
			h.Div(
				h.Class("app-gallery-content__items"),
				g.Group(g.Map(icons, func(i Item) g.Node {
					for _, role := range i.Role {
						if utils.CheckRole(p.Ctx.User.Roles, role) {
							return h.A(
								h.Class("app-gallery-content__item"),
								h.Href(i.Link),
								components.Icon(&components.IconProps{
									Identifier: i.Icon,
								}),
								h.Div(
									h.Class("app-gallery-content__item-name"),
									g.Text(i.Name),
								),
							)
						}
					}
					return g.Text("")
				})),
			),
		),
		components.InlineStyle("/layout/appGallery.css"),
		components.InlineScript("/layout/appGallery.js"),
	)
}
