package components

import (
	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
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
