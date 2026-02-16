package imagetotextview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type ImageToTextFixPageProps struct {
	Ctx reqcontext.ReqContext
}

func ImageToTextFixPage(p *ImageToTextFixPageProps) g.Node {
	content := h.Div(
		h.Class("image-to-text-fix-page"),
		h.Div(
			h.Class("image-to-text-fix-hero"),
			h.H1(g.Text("Fix OCR")),
			h.P(g.Text("Review the extracted text, adjust the regex if needed, and apply the result back to your form.")),
		),
		components.Card(
			c.Classes{"image-to-text-fix-card": true},
			h.Div(
				h.Class("image-to-text-fix-field"),
				h.Label(h.For("ocr-fix-text"), g.Text("Extracted text")),
				h.Textarea(
					h.ID("ocr-fix-text"),
					h.Rows("10"),
				),
			),
			h.Div(
				h.Class("image-to-text-fix-field"),
				h.Label(h.For("ocr-fix-pattern"), g.Text("Regex pattern")),
				h.Textarea(
					h.ID("ocr-fix-pattern"),
					h.Rows("4"),
				),
			),
			h.Div(
				h.Class("image-to-text-fix-field"),
				h.Label(h.For("ocr-fix-flags"), g.Text("Regex flags (optional)")),
				h.Input(
					h.Type("text"),
					h.ID("ocr-fix-flags"),
					h.Placeholder("i, m, s"),
					h.AutoComplete("off"),
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
			),
			h.Div(
				h.Class("image-to-text-fix-error"),
				h.Span(h.ID("ocr-fix-error")),
			),
		),
	)

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Fix OCR",
		Content: content,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/imagetotextview/image_to_text_fix_page.css"),
			components.InlineScript("/internal/components/OcrClient.js"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/imagetotextview/image_to_text_fix_page.js"),
		},
	})
}
