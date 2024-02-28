package layout

import (
	reqContext "app/reqcontext"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type layoutProps struct {
	content   g.Node
	Ctx       reqContext.ReqContext
	NoPadding bool
}

func layout(p *layoutProps) g.Node {
	return g.Group([]g.Node{
		navbar(&navbarProps{
			Ctx: p.Ctx,
		}),
		g.If(p.Ctx.User.UserID != 0,
			breadcrumbs(p.Ctx.Req.URL.Path),
		),
		h.Main(
			c.Classes{
				"main":         true,
				"main-padding": !p.NoPadding,
			},
			p.content,
		),
		footer(),
	})
}
