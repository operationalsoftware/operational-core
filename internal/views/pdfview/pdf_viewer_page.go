package pdfview

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type PDFInlineViewerProps struct {
	Title string
	Src   string
}

// PDFInlineViewer renders a minimal full-page iframe for a PDF URL.
func PDFInlineViewer(p PDFInlineViewerProps) g.Node {
	title := p.Title
	if title == "" {
		title = "PDF Document"
	}

	return h.Doctype(
		h.HTML(
			h.Lang("en"),
			h.Head(
				h.Meta(h.Charset("utf-8")),
				h.TitleEl(g.Text(title)),
				h.StyleEl(g.Text(`
html, body, iframe { margin: 0; padding: 0; height: 100%; width: 100%; }
iframe { border: 0; }
`)),
			),
			h.Body(
				h.IFrame(
					h.Src(p.Src),
					g.Attr("allow", "fullscreen"),
				),
			),
		),
	)
}
