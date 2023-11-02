package layout

import (
	o "operationalcore/components"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type layoutParams struct {
	crumbs  []Crumb
	content g.Node
}

func layout(params layoutParams) []g.Node {
	return []g.Node{
		o.InlineStyle(Assets, "/layout.css"),
		o.InlineScript(Assets, "/global.js"),
		navbar(),
		breadcrumbs(params.crumbs),
		Main(Class("main"), params.content),
		footer(),
	}
}
