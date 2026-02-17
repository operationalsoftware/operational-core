package components

import (
	"encoding/json"
	"fmt"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type OCRPatterns struct {
	Pattern string `json:"pattern"`
	Example string `json:"example"`
	Flags   string `json:"flags"`
}

type OCRInputProps struct {
	Target  string
	Classes c.Classes
	Regex   []OCRPatterns
}

func OCRInput(p *OCRInputProps, children ...g.Node) g.Node {
	if p.Classes == nil {
		p.Classes = c.Classes{}
	}

	p.Classes["ocr-input"] = true

	fileID := fmt.Sprintf("%s-ocr-file", p.Target)
	regexJSON := ""
	if len(p.Regex) > 0 {
		if payload, err := json.Marshal(p.Regex); err == nil {
			regexJSON = string(payload)
		}
	}

	return h.Div(
		p.Classes,
		g.Attr("data-ocr-input", ""),
		g.If(regexJSON != "", g.Attr("data-ocr-regex-list", regexJSON)),
		g.Attr("data-ocr-param", p.Target),
		g.Attr("data-ocr-name", p.Target),
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
