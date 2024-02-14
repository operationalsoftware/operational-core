package layout

import (
	"app/components"
	"app/utils"

	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type item struct {
	icon string
	name string
	link string
}

type appGalleryProps struct {
	Ctx utils.Context
}

func appGallery(p *appGalleryProps) g.Node {
	classes := c.Classes{
		"app-gallery": true,
	}

	// icons array
	icons := []item{
		{
			icon: "account-group",
			name: "Users",
			link: "/users",
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
				g.Group(g.Map(icons, func(i item) g.Node {
					return h.A(
						h.Class("app-gallery-content__item"),
						h.Href(i.link),
						components.Icon(&components.IconProps{
							Identifier: i.icon,
						}),
						h.Div(
							h.Class("app-gallery-content__item-name"),
							g.Text(i.name),
						),
					)
				})),
			),
		),
		components.InlineStyle("/layout/appGallery.css"),
		components.InlineScript("/layout/appGallery.js"),
	)
}
