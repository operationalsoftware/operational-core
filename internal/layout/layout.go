package layout

import (
	"app/pkg/reqcontext"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type layoutProps struct {
	breadcrumbs []Breadcrumb
	content     g.Node
	ctx         reqcontext.ReqContext
	noPadding   bool
}

func layout(p *layoutProps) g.Node {
	return g.Group([]g.Node{
		navbar(&navbarProps{
			ctx: p.ctx,
		}),
		g.If(p.ctx.User.UserID != 0,
			breadcrumbs(p.breadcrumbs),
		),
		h.Main(
			c.Classes{
				"main":         true,
				"main-padding": !p.noPadding,
			},
			p.content,
		),
		footer(),
	})
}
