package service

import (
	"app/pkg/printnode"
	"context"
	"encoding/base64"
	"time"
)

type PrintNodeService struct {
	client *printnode.Client
}

func NewPrintNodeService(apiKey string) *PrintNodeService {
	return &PrintNodeService{
		client: printnode.NewClient(apiKey),
	}
}

func (s *PrintNodeService) Status(ctx context.Context) (printnode.Status, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.client.CheckStatus(ctx)
}

func (s *PrintNodeService) Printers(ctx context.Context) ([]printnode.Printer, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.client.ListPrinters(ctx)
}

func (s *PrintNodeService) SubmitPDF(ctx context.Context, printerID int, title string, pdfBytes []byte) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	encoded := base64.StdEncoding.EncodeToString(pdfBytes)
	return s.client.SubmitPrintJob(ctx, printerID, title, "pdf_base64", encoded)
}
