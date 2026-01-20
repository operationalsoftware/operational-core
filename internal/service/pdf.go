package service

import (
	"app/internal/model"
	"app/internal/pdftemplate"
	"app/internal/repository"
	"app/pkg/pdf"
	"app/pkg/printnode"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ncw/swift/v2"
)

type PDFService struct {
	db        *pgxpool.Pool
	swiftConn *swift.Connection
	fileRepo  *repository.FileRepository
	pdfRepo   *repository.PDFRepository
	printNode *PrintNodeService
}

func NewPDFService(
	db *pgxpool.Pool,
	swiftConn *swift.Connection,
	fileRepo *repository.FileRepository,
	pdfRepo *repository.PDFRepository,
	printNode *PrintNodeService,
) *PDFService {
	return &PDFService{
		db:        db,
		swiftConn: swiftConn,
		fileRepo:  fileRepo,
		pdfRepo:   pdfRepo,
		printNode: printNode,
	}
}

// GenerateFromJSON builds a PDF and returns both the bytes and the resolved title.
func (s *PDFService) GenerateFromJSON(
	ctx context.Context,
	templateName string,
	jsonInput []byte,
) ([]byte, string, error) {
	template, ok := pdftemplate.Registry[templateName]
	if !ok {
		return nil, "", fmt.Errorf("unknown template: %s", templateName)
	}

	pdfDefinition, err := template.Generator.GenerateFromJSON(jsonInput)
	if err != nil {
		return []byte{}, "", fmt.Errorf("error generating PDF definition from template: %v", err)
	}

	title := strings.TrimSpace(pdfDefinition.Title)
	if title == "" {
		title = pdftemplate.FallbackTitle(templateName)
	}
	pdfDefinition.Title = title

	pdfBuffer, err := pdf.GeneratePDF(ctx, pdfDefinition)
	if err != nil {
		return []byte{}, "", fmt.Errorf("error generating PDF from definition: %v", err)
	}

	return pdfBuffer, title, nil
}

// RecordGeneration persists the generated PDF and a log entry. It returns the log record.
func (s *PDFService) RecordGeneration(
	ctx context.Context,
	templateName string,
	inputData string,
	pdfBytes []byte,
	userID int,
	pdfTitle string,
) (model.PDFGenerationLog, error) {
	if s.db == nil || s.swiftConn == nil || s.fileRepo == nil || s.pdfRepo == nil {
		return model.PDFGenerationLog{}, fmt.Errorf("pdf service not fully configured for logging")
	}
	if userID == 0 {
		return model.PDFGenerationLog{}, fmt.Errorf("user id is required for logging")
	}
	title := strings.TrimSpace(pdfTitle)
	if title == "" {
		title = pdftemplate.FallbackTitle(templateName)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return model.PDFGenerationLog{}, err
	}
	defer tx.Rollback(ctx)

	logID, err := s.pdfRepo.InsertGenerationLog(ctx, tx, templateName, inputData, userID, title)
	if err != nil {
		return model.PDFGenerationLog{}, fmt.Errorf("failed to insert pdf generation log: %w", err)
	}

	fileName := s.GeneratePDFFilename(title)

	file, err := s.fileRepo.SaveFileContent(ctx, tx, s.swiftConn, &model.File{
		Filename:    fileName,
		ContentType: "application/pdf",
		SizeBytes:   len(pdfBytes),
		Entity:      "PDFGenerationLog",
		EntityID:    logID,
		UserID:      userID,
	}, pdfBytes)
	if err != nil {
		return model.PDFGenerationLog{}, fmt.Errorf("failed to store generated pdf: %w", err)
	}

	if err := s.pdfRepo.UpdateGenerationLogFile(ctx, tx, logID, file.FileID, title); err != nil {
		return model.PDFGenerationLog{}, fmt.Errorf("failed to update pdf generation log with file: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.PDFGenerationLog{}, err
	}

	downloadURL, err := s.fileRepo.GetSignedDownloadURLWithDisposition(
		ctx,
		s.swiftConn,
		s.db,
		file.FileID,
		15*time.Minute,
		file.Filename,
		true,
	)
	if err != nil {
		return model.PDFGenerationLog{}, fmt.Errorf("failed to generate download url: %w", err)
	}

	return model.PDFGenerationLog{
		PDFGenerationLogID: logID,
		TemplateName:       templateName,
		InputData:          inputData,
		FileID:             file.FileID,
		PDFTitle:           title,
		FileURL:            downloadURL,
		UserID:             userID,
		CreatedAt:          time.Now(),
	}, nil
}

func (s *PDFService) ListRecentLogs(ctx context.Context, limit int) ([]model.PDFGenerationLog, error) {
	if s.db == nil || s.pdfRepo == nil {
		return nil, fmt.Errorf("pdf service not configured for listing logs")
	}

	logs, err := s.pdfRepo.ListRecentGenerationLogs(ctx, s.db, limit)
	if err != nil {
		return nil, err
	}

	for i := range logs {
		if logs[i].FileID == "" {
			continue
		}
		url, err := s.fileRepo.GetSignedDownloadURLWithDisposition(
			ctx,
			s.swiftConn,
			s.db,
			logs[i].FileID,
			15*time.Minute,
			s.GeneratePDFFilename(logs[i].PDFTitle),
			true,
		)
		if err != nil {
			continue
		}
		logs[i].FileURL = url
	}

	return logs, nil
}

func (s *PDFService) ListGenerationLogs(ctx context.Context, limit, offset int) ([]model.PDFGenerationLog, int, error) {
	if s.db == nil || s.pdfRepo == nil {
		return nil, 0, fmt.Errorf("pdf service not configured for listing logs")
	}

	total, err := s.pdfRepo.CountGenerationLogs(ctx, s.db)
	if err != nil {
		return nil, 0, err
	}

	logs, err := s.pdfRepo.ListGenerationLogs(ctx, s.db, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	for i := range logs {
		if logs[i].FileID == "" {
			continue
		}
		url, err := s.fileRepo.GetSignedDownloadURLWithDisposition(
			ctx,
			s.swiftConn,
			s.db,
			logs[i].FileID,
			15*time.Minute,
			s.GeneratePDFFilename(logs[i].PDFTitle),
			true,
		)
		if err != nil {
			continue
		}
		logs[i].FileURL = url
	}

	return logs, total, nil
}

func (s *PDFService) fetchGenerationLog(ctx context.Context, id int) (model.PDFGenerationLog, error) {
	return s.pdfRepo.GetGenerationLog(ctx, s.db, id)
}

func (s *PDFService) PrintAndLog(
	ctx context.Context,
	templateName string,
	inputData string,
	printerName string,
	requirementName string,
	userID int,
) (model.PDFPrintLog, error) {
	if s.printNode == nil {
		return model.PDFPrintLog{}, fmt.Errorf("printnode service not configured")
	}
	printerName = strings.TrimSpace(printerName)
	if printerName == "" {
		return model.PDFPrintLog{}, fmt.Errorf("printer name is required")
	}

	pdfBytes, title, err := s.GenerateFromJSON(ctx, templateName, []byte(inputData))
	if err != nil {
		return model.PDFPrintLog{}, err
	}

	genLog, err := s.RecordGeneration(ctx, templateName, inputData, pdfBytes, userID, title)
	if err != nil {
		return model.PDFPrintLog{}, err
	}

	jobID := 0
	status := "success"
	errorMessage := ""

	printerID := 0
	printers, err := s.printNode.Printers(ctx)
	if err != nil {
		status = "error"
		errorMessage = err.Error()
	}

	for _, pr := range printers {
		if strings.EqualFold(pr.Name, printerName) {
			printerID = pr.ID
			break
		}
	}
	if printerID == 0 {
		err = fmt.Errorf("printer not found: %s", printerName)
		status = "error"
		errorMessage = err.Error()
	}

	jobID, err = s.printNode.SubmitPDF(ctx, printerID, templateName, pdfBytes)
	if err != nil {
		status = "error"
		errorMessage = err.Error()
	}

	printLogID := 0
	if s.db != nil && s.pdfRepo != nil {
		id, insertErr := s.pdfRepo.InsertPrintLog(ctx, s.db, genLog.PDFGenerationLogID, templateName, requirementName, printerName, jobID, status, errorMessage, userID)
		if insertErr == nil {
			printLogID = id
		}
	}

	return model.PDFPrintLog{
		PDFPrintLogID:      printLogID,
		PDFGenerationLogID: genLog.PDFGenerationLogID,
		TemplateName:       templateName,
		InputData:          inputData,
		RequirementName:    requirementName,
		PrinterName:        printerName,
		PrintNodeJobID:     jobID,
		Status:             status,
		ErrorMessage:       errorMessage,
		FileID:             genLog.FileID,
		PDFTitle:           genLog.PDFTitle,
		FileURL:            genLog.FileURL,
		UserID:             userID,
		CreatedAt:          time.Now(),
	}, err
}

func (s *PDFService) Reprint(ctx context.Context, printLogID int, overridePrinterName string, userID int) (model.PDFPrintLog, error) {
	var existing model.PDFPrintLog
	existing, err := s.pdfRepo.GetPrintLog(ctx, s.db, printLogID)
	if err != nil {
		return model.PDFPrintLog{}, err
	}

	logEntry, err := s.fetchGenerationLog(ctx, existing.PDFGenerationLogID)
	if err != nil {
		return model.PDFPrintLog{}, err
	}

	printerName := existing.PrinterName
	if strings.TrimSpace(overridePrinterName) != "" {
		printerName = overridePrinterName
	}

	return s.PrintAndLog(ctx, logEntry.TemplateName, logEntry.InputData, printerName, existing.RequirementName, userID)
}

func (s *PDFService) ListRecentPrintLogs(ctx context.Context, limit int) ([]model.PDFPrintLog, error) {
	if s.db == nil || s.pdfRepo == nil {
		return nil, fmt.Errorf("pdf service not configured for listing print logs")
	}

	logs, err := s.pdfRepo.ListRecentPrintLogs(ctx, s.db, limit)
	if err != nil {
		return nil, err
	}

	for i := range logs {
		if logs[i].FileID == "" {
			continue
		}
		url, err := s.fileRepo.GetSignedDownloadURLWithDisposition(
			ctx,
			s.swiftConn,
			s.db,
			logs[i].FileID,
			15*time.Minute,
			s.GeneratePDFFilename(logs[i].PDFTitle),
			true,
		)
		if err != nil {
			continue
		}
		logs[i].FileURL = url
	}

	return logs, nil
}

func (s *PDFService) ListPrintRequirements(ctx context.Context) ([]model.PrintRequirement, error) {
	if s.db == nil || s.pdfRepo == nil {
		return nil, fmt.Errorf("pdf service not configured for listing requirements")
	}
	return s.pdfRepo.ListPrintRequirements(ctx, s.db)
}

func (s *PDFService) SavePrintRequirement(ctx context.Context, pr model.PrintRequirement) (model.PrintRequirement, error) {
	if s.db == nil || s.pdfRepo == nil {
		return model.PrintRequirement{}, fmt.Errorf("pdf service not configured for saving requirements")
	}
	pr.RequirementName = strings.TrimSpace(pr.RequirementName)
	if pr.RequirementName == "" {
		return model.PrintRequirement{}, fmt.Errorf("requirement name is required")
	}
	pr.PrinterName = strings.TrimSpace(pr.PrinterName)
	if pr.PrinterName == "" {
		return model.PrintRequirement{}, fmt.Errorf("printer name is required")
	}

	assignments, err := s.pdfRepo.ListPrintRequirements(ctx, s.db)
	if err != nil {
		return model.PrintRequirement{}, err
	}

	for _, a := range assignments {
		if pr.PrinterName == "" {
			break
		}
		if strings.EqualFold(a.PrinterName, pr.PrinterName) && !strings.EqualFold(a.RequirementName, pr.RequirementName) {
			return model.PrintRequirement{}, fmt.Errorf("printer already assigned to another requirement")
		}
	}

	return s.pdfRepo.UpsertPrintRequirement(ctx, s.db, pr)
}

func (s *PDFService) ListAvailablePrinters(ctx context.Context, currentReq string, printers []printnode.Printer) ([]printnode.Printer, error) {
	if s.db == nil || s.pdfRepo == nil {
		return nil, fmt.Errorf("pdf service not configured for listing requirements")
	}
	assignments, err := s.pdfRepo.ListPrintRequirements(ctx, s.db)
	if err != nil {
		return nil, err
	}
	assigned := map[string]struct{}{}
	for _, a := range assignments {
		if strings.EqualFold(a.RequirementName, currentReq) {
			continue
		}
		if name := strings.TrimSpace(a.PrinterName); name != "" {
			assigned[strings.ToLower(name)] = struct{}{}
		}
	}
	available := make([]printnode.Printer, 0, len(printers))
	for _, p := range printers {
		if _, taken := assigned[strings.ToLower(strings.TrimSpace(p.Name))]; taken {
			continue
		}
		available = append(available, p)
	}
	return available, nil
}

func (s *PDFService) GeneratePDFTitleFromInput(templateName string, jsonInput []byte) string {
	return pdftemplate.FallbackTitle(templateName)
}

// GeneratePDFFilename returns a filename (with .pdf) based on the PDF title.
func (s *PDFService) GeneratePDFFilename(pdfTitle string) string {
	name := strings.TrimSpace(pdfTitle)
	if !strings.HasSuffix(strings.ToLower(name), ".pdf") {
		name += ".pdf"
	}
	return name
}
