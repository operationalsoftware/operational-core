package service

import (
	"app/internal/pdftemplate"
	"app/pkg/pdf"
	"context"
	"fmt"
)

type PDFService struct {}

func NewPDFService() *PDFService {
	return &PDFService{}
}

func (s *PDFService) GenerateFromJSON(
	ctx context.Context,
	templateName string,
	jsonInput []byte,
) ([]byte, error) {
	template, ok := pdftemplate.Registry[templateName]
	if !ok {
		return nil, fmt.Errorf("unknown template: %s", templateName)
	}

	pdfDefinition, err := template.Generator.GenerateFromJSON(jsonInput)
	if err != nil {
		return []byte{}, fmt.Errorf("error generating PDF definition from template: %v", err)
	}

	pdfBuffer, err := pdf.GeneratePDF(ctx, pdfDefinition)
	if err != nil {
		return []byte{}, fmt.Errorf("error generating PDF from definition: %v", err)
	}

	return pdfBuffer, nil
}
