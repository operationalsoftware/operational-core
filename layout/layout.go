package layout

import (
	o "operationalcore/components"
	"operationalcore/utils"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type layoutProps struct {
	crumbs    []Crumb
	content   g.Node
	Ctx       utils.Context
	NoPadding bool
}

func layout(p *layoutProps) []g.Node {
	return []g.Node{
		o.InlineStyle(Assets, "/layout.css"),
		o.InlineScript(Assets, "/global.js"),
		navbar(&navbarProps{
			Ctx: p.Ctx,
		}),
		breadcrumbs(p.crumbs),
		h.Main(
			c.Classes{
				"main":         true,
				"main-padding": !p.NoPadding,
			},
			p.content,
		),
		footer(),
	}
}
