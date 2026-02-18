package imagetotextview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type ImageToTextResolvePageProps struct {
	Ctx reqcontext.ReqContext
}

func ImageToTextResolvePage(p *ImageToTextResolvePageProps) g.Node {
	content := h.Div(
		h.Class("image-to-text-fix-page"),
		h.Div(
			h.Class("image-to-text-fix-hero"),
			h.H1(g.Text("Resolve OCR")),
			h.P(g.Text("Review the extracted text and apply the result back to your form.")),
		),
		components.Card(
			c.Classes{"image-to-text-fix-card": true},
			h.Div(
				h.Class("image-to-text-fix-field"),
				h.Label(g.Text("Extracted text tokens")),
				h.Div(
					h.ID("ocr-fix-tags"),
					h.Class("image-to-text-fix-tags"),
				),
				h.Label(h.For("ocr-fix-input"), g.Text("Selected text")),
				h.Input(
					h.Type("text"),
					h.ID("ocr-fix-input"),
					h.Placeholder("Click tokens or type to build the value"),
					h.AutoComplete("off"),
				),
				h.P(
					h.Class("image-to-text-fix-example"),
					h.ID("ocr-fix-example"),
				),
			),
			h.Div(
				h.Class("image-to-text-fix-status"),
				h.Span(h.ID("ocr-fix-status"), g.Text("Ready.")),
			),
			h.Div(
				h.Class("image-to-text-fix-actions"),
				h.Button(
					h.Type("button"),
					h.Class("button primary"),
					h.ID("ocr-fix-submit"),
					g.Text("Apply to form"),
				),
				h.A(
					h.Class("button secondary"),
					h.ID("ocr-fix-back"),
					h.Href("#"),
					g.Text("Go back"),
				),
			),
			h.Div(
				h.Class("image-to-text-fix-error"),
				h.Span(h.ID("ocr-fix-error")),
			),
		),
	)

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Resolve OCR",
		Content: content,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/imagetotextview/resolve.css"),
			h.Script(h.Src("/static/js/lib/ocr-client.js")),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/imagetotextview/resolve.js"),
		},
	})
}
