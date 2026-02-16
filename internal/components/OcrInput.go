package components

import (
	"fmt"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type OcrInputProps struct {
	Size         InputSize
	Name         string
	Label        string
	Placeholder  string
	HelperText   string
	InputType    string
	HelperType   InputHelperType
	InputProps   []g.Node
	Classes      c.Classes
	RegexPattern string
	RegexFlags   string
	ParamName    string
}

func OcrInput(p *OcrInputProps, children ...g.Node) g.Node {
	if p.Classes == nil {
		p.Classes = c.Classes{}
	}

	if p.InputType == "" {
		p.InputType = "text"
	}

	if p.InputProps == nil {
		p.InputProps = []g.Node{}
	}

	if p.Size == "" {
		p.Size = InputSizeMedium
	}

	if p.ParamName == "" {
		p.ParamName = p.Name
	}

	inputClasses := c.Classes{}
	inputClasses[string(p.Size)] = true

	p.Classes["input-container"] = true
	p.Classes["ocr-input"] = true

	fileID := fmt.Sprintf("%s-ocr-file", p.Name)

	return h.Div(
		p.Classes,
		g.Attr("data-ocr-input", ""),
		g.Attr("data-ocr-pattern", p.RegexPattern),
		g.Attr("data-ocr-flags", p.RegexFlags),
		g.Attr("data-ocr-param", p.ParamName),
		g.If(
			p.Label != "",
			h.Label(h.For(p.Name), g.Text(p.Label)),
		),
		h.Div(
			h.Class("ocr-input-row"),
			h.Input(
				inputClasses,
				h.Name(p.Name),
				h.ID(p.Name),
				h.Placeholder(p.Placeholder),
				h.Type(p.InputType),
				g.Attr("data-ocr-field", ""),
				g.Group(p.InputProps),
			),
			h.Button(
				h.Type("button"),
				h.Class("button ocr-input-button"),
				g.Attr("data-ocr-trigger", ""),
				g.Attr("aria-label", "Scan with camera"),
				Icon(&IconProps{Identifier: "camera"}),
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
		h.Div(
			h.Class("ocr-input-meta"),
			h.Span(h.Class("ocr-input-status"), g.Attr("data-ocr-status", ""), g.Text("Ready.")),
			h.Span(h.Class("ocr-input-error"), g.Attr("data-ocr-error", "")),
		),
		g.If(p.HelperText != "", InputHelper(&InputHelperProps{
			Label: p.HelperText,
			Type:  p.HelperType,
		})),
		g.Group(children),
		InlineStyle("/internal/components/OcrInput.css"),
		h.Script(h.Src("https://cdn.jsdelivr.net/npm/tesseract.js@5/dist/tesseract.min.js"), g.Attr("defer", "")),
		InlineScript("/internal/components/OcrClient.js"),
		InlineScript("/internal/components/OcrInput.js"),
	)
}
