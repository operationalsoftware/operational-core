package pdftemplate

import (
	"app/pkg/pdf"
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

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
	Name             string
	Description      string
	Generator        TemplateGenerator
	ExampleJSON      string
	PrintNodeOptions json.RawMessage
}

var Registry = map[string]RegisteredTemplate{
	InvoiceTemplateDefinition.Name: InvoiceTemplateDefinition,
}

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

func FallbackTitle(base string) string {
	title := strings.TrimSpace(base)
	if title == "" {
		title = "PDF"
	}
	return fmt.Sprintf("%s-%s", title, time.Now().Format("200601021504"))
}
