package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

func CopyButton(textToCopy string) g.Node {

	return h.Div(
		h.Class("clipboard-container"),
		h.DataAttr("text", textToCopy),
		h.Button(
			h.Class("clipboard-btn"),
			Icon(&IconProps{
				Identifier: "content-copy",
				Classes: c.Classes{
					"icon": true,
				},
			}),
		),

		h.Span(
			h.Class("clipboard-status hidden"),
			g.Text("Copied!"),
		),

		InlineScript("/internal/components/CopyButton.js"),
	)
}
