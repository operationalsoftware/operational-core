package layout

import (
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func footer() g.Node {
	return Footer(
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
