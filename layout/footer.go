package layout

import (
	o "operationalcore/components"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func footer() g.Node {
	return Footer(
		o.InlineStyle(
			Assets, "/footer.css",
		),
		Div(
			g.Text("An OperationalPlatform"),
			Sup(g.Text("TM")),
			g.Text(" by "),
			A(
				Href("https://operationalsoftware.co"),
				Target("_blank"),
				g.Text("Operational Software"),
			),
		),
	)
}
