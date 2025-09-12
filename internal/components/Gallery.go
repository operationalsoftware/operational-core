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
			return h.Div(
				h.Class("gallery-item"),
				h.Img(
					h.Src(src),
					h.Alt("gallery image"),
				),
			)
		})),

		InlineStyle("/internal/components/Gallery.css"),
		InlineScript("/internal/components/Gallery.js"),
	)
}
