package pdftemplate

import (
	"app/pkg/pdf"
	"bytes"
	"encoding/json"
	"sort"

	g "maragu.dev/gomponents"
)

// utility to convert a gomponent to a string
func gomponentToString(node g.Node) (string, error) {
	var buf bytes.Buffer
	if err := node.Render(&buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func GenerateTypedFromJSON[T any](typedGenerate func(T) (pdf.PDFDefinition, error), jsonInput []byte) (pdf.PDFDefinition, error) {
	var input T
	if err := json.Unmarshal(jsonInput, &input); err != nil {
		return pdf.PDFDefinition{}, err
	}
	return typedGenerate(input)
}

type TemplateGenerator interface {
	GenerateFromJSON([]byte) (pdf.PDFDefinition, error)
}

type RegisteredTemplate struct {
	Name        string
	Description string
	Generator   TemplateGenerator
	ExampleJSON string
	// TitleGenerator is optional and can derive a PDF title from the JSON input.
	TitleGenerator func(jsonInput []byte) (string, error)
}

var Registry = func() map[string]RegisteredTemplate {
	return map[string]RegisteredTemplate{
		InvoiceTemplateDefinition.Name: InvoiceTemplateDefinition,
	}
}()

// SortedTemplates returns a slice of RegisteredTemplate sorted by Name.
func SortedTemplates() []RegisteredTemplate {
	templates := make([]RegisteredTemplate, 0, len(Registry))
	for _, tmpl := range Registry {
		templates = append(templates, tmpl)
	}

	sort.Slice(templates, func(i, j int) bool {
		return templates[i].Name < templates[j].Name
	})

	return templates
}
