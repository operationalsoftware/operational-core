package components

import (
	"fmt"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type OCRInputProps struct {
	Size         InputSize
	Name         string
	Placeholder  string
	Classes      c.Classes
	RegexPattern string
	RegexFlags   string
	ParamName    string
}

func OCRInput(p *OCRInputProps, children ...g.Node) g.Node {
	if p.Classes == nil {
		p.Classes = c.Classes{}
	}

	if p.Size == "" {
		p.Size = InputSizeMedium
	}

	if p.ParamName == "" {
		p.ParamName = p.Name
	}

	p.Classes["ocr-input"] = true

	fileID := fmt.Sprintf("%s-ocr-file", p.Name)

	return h.Div(
		p.Classes,
		g.Attr("data-ocr-input", ""),
		g.Attr("data-ocr-pattern", p.RegexPattern),
		g.Attr("data-ocr-flags", p.RegexFlags),
		g.Attr("data-ocr-param", p.ParamName),
		g.Attr("data-ocr-name", p.Name),
		g.If(p.Placeholder != "", g.Attr("data-ocr-example", p.Placeholder)),
		h.Div(
			h.Class("ocr-input-row"),
			h.Button(
				h.Type("button"),
				h.Class("button ocr-input-button"),
				g.Attr("data-ocr-trigger", ""),
				g.Attr("aria-label", "Scan with OCR"),
				Icon(&IconProps{Identifier: "ocr"}),
			),
			h.Input(
				h.Type("file"),
				h.ID(fileID),
				h.Class("ocr-input-file"),
				g.Attr("hidden", ""),
				g.Attr("accept", "image/*"),
				g.Attr("capture", "environment"),
				g.Attr("data-ocr-file", ""),
			),
		),
		g.Group(children),
		InlineStyle("/internal/components/OCRInput.css"),
		h.Script(h.Src("https://cdn.jsdelivr.net/npm/tesseract.js@5/dist/tesseract.min.js"), g.Attr("defer", "")),
		h.Script(h.Src("/static/js/lib/ocr-client.js")),
		InlineScript("/internal/components/OCRInput.js"),
	)
}
