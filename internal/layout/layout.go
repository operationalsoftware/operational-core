package layout

import (
	"app/internal/components"
	"app/pkg/reqcontext"
	"os"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type layoutProps struct {
	breadcrumbs []Breadcrumb
	content     g.Node
	ctx         reqcontext.ReqContext
	mainPadding *bool
}

func layout(p *layoutProps) g.Node {

	appEnv := os.Getenv("APP_ENV")
	showBanner := appEnv == "staging"

	mainPadding := true
	if p.mainPadding != nil {
		mainPadding = *p.mainPadding
	}

	return g.Group([]g.Node{
		h.Button(
			h.ID("navbar-expand-menu-button"),
			h.Class("menu-button navbar-expand-button hidden"),
			components.Icon(&components.IconProps{
				Identifier: "chevron-double-down",
			}),
		),
		g.If(
			showBanner,
			h.Div(
				h.Class("env-banner"),

				g.Text("Staging"),
			),
		),
		navbar(&navbarProps{
			ctx: p.ctx,
		}),
		g.If(p.ctx.User.UserID != 0,
			breadcrumbs(p.breadcrumbs),
		),
		h.Main(
			c.Classes{
				"main":         true,
				"main-padding": mainPadding,
			},
			p.content,
		),
		footer(),
	})
}
