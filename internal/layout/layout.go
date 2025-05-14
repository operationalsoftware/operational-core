package layout

import (
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
	showBanner := appEnv != "dev" && appEnv != "production"

	mainPadding := true
	if p.mainPadding != nil {
		mainPadding = *p.mainPadding
	}

	return g.Group([]g.Node{
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
