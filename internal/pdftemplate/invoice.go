package pdftemplate

import (
	"app/pkg/pdf"
	"fmt"

	"github.com/shopspring/decimal"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type InvoiceData struct {
	CustomerName string
	Amount       decimal.Decimal
}

type InvoiceTemplate struct{}

func (InvoiceTemplate) Generate(input InvoiceData) (pdf.PDFDefinition, error) {

	html, err := gomponentToString(h.Div(
		h.H1(g.Textf("Invoice for %s", input.CustomerName)),
		h.P(g.Textf("Amount: %s", input.Amount.String())),
	))
	if err != nil {
		return pdf.PDFDefinition{}, fmt.Errorf("error generating invoice html: %v", err)
	}

	title := "Invoice"

	return pdf.PDFDefinition{Title: title, HTML: html}, nil
}

func (InvoiceTemplate) GenerateFromJSON(data []byte) (pdf.PDFDefinition, error) {
	return GenerateTypedFromJSON(InvoiceTemplate{}.Generate, data)
}

var invoiceExampleJSON = `
{
  "CustomerName": "Jane Doe",
  "Amount": 123.45
}`

var InvoiceTemplateDefinition = RegisteredTemplate{
	Name:        "Invoice",
	Description: "Simple invoice with customer name and amount",
	Generator:   InvoiceTemplate{},
	ExampleJSON: invoiceExampleJSON,
}
