package service

import (
	"app/internal/pdftemplate"
	"app/internal/repository"
	"app/pkg/pdf"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PDFService struct {
	db            *pgxpool.Pool
	pdfRepository *repository.PDFRepository
}

func NewPDFService(
	db *pgxpool.Pool,
	pdfRepository *repository.PDFRepository,
) *PDFService {
	return &PDFService{
		db:            db,
		pdfRepository: pdfRepository,
	}
}

func (s *PDFService) GeneratePDF(
	ctx context.Context,
	template string,
	inputParams map[string]interface{},
) ([]byte, error) {
	tmpl, ok := pdftemplate.TemplateRegistry[template]
	if !ok {
		return nil, errors.New("unknown template: " + template)
	}

	pdfOptions := pdf.PdfOptions{Title: "My Report"}

	pdfBuffer, err := pdf.GeneratePDF(tmpl(nil).HTML, &pdfOptions)
	if err != nil {
		return nil, errors.New("failed to generate pdf")
	}

	return pdfBuffer, nil
}
