package components

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func CopyButton(textToCopy string) g.Node {

	return h.Div(
		h.Class("copy-button"),
		h.DataAttr("text", textToCopy),
		h.Button(
			h.Class("button"),
			Icon(&IconProps{Identifier: "content-copy"}),
		),

		h.Span(
			h.Class("status hidden"),
			g.Text("Copied!"),
		),

		InlineScript("/internal/components/CopyButton.js"),
	)
}
