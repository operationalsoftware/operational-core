package imagetotextview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type ImageToTextPageProps struct {
	Ctx reqcontext.ReqContext
}

func ImageToTextPage(p *ImageToTextPageProps) g.Node {
	content := h.Div(
		h.Class("image-to-text-page"),
		h.Div(
			h.Class("image-to-text-hero"),
			h.H1(g.Text("Image to Text")),
			h.P(g.Text("Upload or capture a document, extract text locally with OCR, and return the captured fields back to your form.")),
		),
		h.Div(
			h.Class("image-to-text-grid"),
			components.Card(
				c.Classes{"image-to-text-card": true},
				h.H2(g.Text("Setup")),
				h.P(
					g.Text("Provide a return URL and a regex pattern with named capture groups. Matches are appended to the return URL as query parameters."),
				),
				h.Div(
					h.Class("image-to-text-field"),
					h.Label(h.For("return-url"), g.Text("Return URL")),
					h.Input(
						h.Type("url"),
						h.ID("return-url"),
						h.Placeholder("https://app.example.com/form?state=123"),
						h.AutoComplete("off"),
					),
					h.P(
						h.Class("image-to-text-hint"),
						g.Text("Query params on this page: return_to (or returnTo), pattern, flags."),
					),
				),
				h.Div(
					h.Class("image-to-text-field"),
					h.Label(h.For("regex-pattern"), g.Text("Regex pattern")),
					h.Textarea(
						h.ID("regex-pattern"),
						h.Placeholder("(?<serial>[A-Z0-9-]+)"),
						h.Rows("4"),
					),
					h.P(
						h.Class("image-to-text-hint"),
						g.Text("Use named capture groups like (?<field>...). Only the first match is used."),
					),
				),
				h.Div(
					h.Class("image-to-text-field"),
					h.Label(h.For("regex-flags"), g.Text("Regex flags (optional)")),
					h.Input(
						h.Type("text"),
						h.ID("regex-flags"),
						h.Placeholder("i, m, s"),
						h.AutoComplete("off"),
					),
				),
			),
			components.Card(
				c.Classes{"image-to-text-card": true},
				h.H2(g.Text("Capture")),
				h.P(g.Text("Choose an image or take a photo. OCR runs on this page; no data is sent to the server.")),
				h.Div(
					h.Class("image-to-text-upload"),
					h.Input(
						h.Type("file"),
						h.ID("ocr-file"),
						h.Accept("image/*"),
						g.Attr("capture", "environment"),
						h.Class("image-to-text-file-input"),
					),
					h.Button(
						h.Type("button"),
						h.Class("button secondary"),
						h.ID("ocr-reset"),
						g.Text("Reset"),
					),
				),
				h.Div(
					h.Class("image-to-text-preview"),
					h.Img(
						h.ID("ocr-preview"),
						h.Class("is-hidden"),
						h.Alt("Selected document preview"),
					),
				),
				h.Div(
					h.Class("image-to-text-status"),
					h.Span(h.ID("ocr-status"), g.Text("Ready.")),
				),
				h.Div(
					h.Class("image-to-text-field"),
					h.Label(h.For("ocr-text"), g.Text("Extracted text")),
					h.Textarea(
						h.ID("ocr-text"),
						h.Rows("8"),
						h.ReadOnly(),
					),
				),
				h.Div(
					h.Class("image-to-text-actions"),
				),
				h.Div(
					h.Class("image-to-text-error"),
					h.Span(h.ID("ocr-error")),
				),
			),
		),
	)

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Image to Text",
		Content: content,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/imagetotextview/image_to_text_page.css"),
			h.Script(h.Src("https://cdn.jsdelivr.net/npm/tesseract.js@5/dist/tesseract.min.js"), g.Attr("defer", "")),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/imagetotextview/image_to_text_page.js"),
		},
	})
}
