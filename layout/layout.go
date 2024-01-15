package layout

import (
	o "operationalcore/components"
	"operationalcore/utils"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type layoutProps struct {
	crumbs  []Crumb
	content g.Node
	Ctx     utils.Context
}

func layout(p *layoutProps) []g.Node {
	return []g.Node{
		o.InlineStyle(Assets, "/layout.css"),
		o.InlineScript(Assets, "/global.js"),
		navbar(&navbarProps{
			Ctx: p.Ctx,
		}),
		breadcrumbs(p.crumbs),
		Main(Class("main"), p.content),
		footer(),
	}
}
