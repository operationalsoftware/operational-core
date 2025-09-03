package components

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Gallery(imageURLs []string) g.Node {
	if len(imageURLs) == 0 {
		return nil
	}

	return h.Div(
		h.Class("gallery"),
		g.Group(g.Map(imageURLs, func(src string) g.Node {
			return h.Img(
				h.Class("gallery-item"),
				h.Src(src),
				h.Alt("attachment"),
			)
		})),

		InlineStyle("/internal/components/Gallery.css"),
		InlineScript("/internal/components/Gallery.js"),
	)
}
