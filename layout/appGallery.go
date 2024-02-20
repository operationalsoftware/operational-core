package layout

import (
	"app/components"
	userModel "app/src/users/model"
	"app/utils"

	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type AppGalleryModule struct {
	Icon string
	Name string
	Link string
	Show func(userModel.UserRoles) bool
}

// icons array
var AppGalleryModules = []AppGalleryModule{
	{
		Icon: "account-group",
		Name: "Users",
		Link: "/users",
		Show: func(roles userModel.UserRoles) bool {
			return roles.UserAdmin.Access
		},
	},
}

type appGalleryProps struct {
	Ctx utils.Context
}

func appGallery(p *appGalleryProps) g.Node {
	classes := c.Classes{
		"app-gallery": true,
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
				g.Group(g.Map(AppGalleryModules, func(i AppGalleryModule) g.Node {

					// check if user has role
					show := i.Show(p.Ctx.User.Roles)
					if !show {
						return g.Text("")
					}

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
				})),
			),
		),
		components.InlineStyle("/layout/appGallery.css"),
		components.InlineScript("/layout/appGallery.js"),
	)
}
