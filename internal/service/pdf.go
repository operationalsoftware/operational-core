package service

import (
	"app/internal/model"
	"app/internal/pdftemplate"
	"app/internal/repository"
	"app/pkg/pdf"
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

func (s *PDFService) GenerateFromJSON(
	ctx context.Context,
	templateName string,
	jsonInput []byte,
	pdfTitle string,
) ([]byte, error) {
	template, ok := pdftemplate.Registry[templateName]
	if !ok {
		return nil, fmt.Errorf("unknown template: %s", templateName)
	}

	pdfDefinition, err := template.Generator.GenerateFromJSON(jsonInput)
	if err != nil {
		return []byte{}, fmt.Errorf("error generating PDF definition from template: %v", err)
	}

	pdfTitle = s.resolvePDFTitle(templateName, pdfTitle, jsonInput)
	pdfDefinition.Title = pdfTitle

	pdfBuffer, err := pdf.GeneratePDF(ctx, pdfDefinition)
	if err != nil {
		return []byte{}, fmt.Errorf("error generating PDF from definition: %v", err)
	}

	return pdfBuffer, nil
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
	pdfTitle = s.resolvePDFTitle(templateName, pdfTitle, []byte(inputData))

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return model.PDFGenerationLog{}, err
	}
	defer tx.Rollback(ctx)

	logID, err := s.pdfRepo.InsertGenerationLog(ctx, tx, templateName, inputData, userID, pdfTitle)
	if err != nil {
		return model.PDFGenerationLog{}, fmt.Errorf("failed to insert pdf generation log: %w", err)
	}

	fileName := buildPDFFilename(pdfTitle)

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

	if err := s.pdfRepo.UpdateGenerationLogFile(ctx, tx, logID, file.FileID, pdfTitle); err != nil {
		return model.PDFGenerationLog{}, fmt.Errorf("failed to update pdf generation log with file: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.PDFGenerationLog{}, err
	}

	downloadURL, err := s.fileRepo.GetSignedDownloadURL(ctx, s.swiftConn, s.db, file.FileID, 15*time.Minute)
	if err != nil {
		return model.PDFGenerationLog{}, fmt.Errorf("failed to generate download url: %w", err)
	}

	return model.PDFGenerationLog{
		ID:           logID,
		TemplateName: templateName,
		InputData:    inputData,
		FileID:       file.FileID,
		PDFTitle:     pdfTitle,
		FileURL:      downloadURL,
		UserID:       userID,
		CreatedAt:    time.Now(),
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
		url, err := s.fileRepo.GetSignedDownloadURL(ctx, s.swiftConn, s.db, logs[i].FileID, 15*time.Minute)
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
		url, err := s.fileRepo.GetSignedDownloadURL(ctx, s.swiftConn, s.db, logs[i].FileID, 15*time.Minute)
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
	printerID int,
	printerName string,
	requirementName string,
	userID int,
) (model.PDFPrintLog, error) {
	if s.printNode == nil {
		return model.PDFPrintLog{}, fmt.Errorf("printnode service not configured")
	}
	if printerID == 0 {
		return model.PDFPrintLog{}, fmt.Errorf("printer id is required")
	}

	pdfTitle := s.resolvePDFTitle(templateName, "", []byte(inputData))

	pdfBytes, err := s.GenerateFromJSON(ctx, templateName, []byte(inputData), pdfTitle)
	if err != nil {
		return model.PDFPrintLog{}, err
	}

	genLog, err := s.RecordGeneration(ctx, templateName, inputData, pdfBytes, userID, pdfTitle)
	if err != nil {
		return model.PDFPrintLog{}, err
	}

	jobID, err := s.printNode.SubmitPDF(ctx, printerID, templateName, pdfBytes)
	status := "success"
	errorMessage := ""
	if err != nil {
		status = "error"
		errorMessage = err.Error()
	}

	printLogID := 0
	if s.db != nil && s.pdfRepo != nil {
		id, insertErr := s.pdfRepo.InsertPrintLog(ctx, s.db, genLog.ID, templateName, requirementName, printerID, printerName, jobID, status, errorMessage, userID)
		if insertErr == nil {
			printLogID = id
		}
	}

	return model.PDFPrintLog{
		ID:                 printLogID,
		PDFGenerationLogID: genLog.ID,
		TemplateName:       templateName,
		InputData:          inputData,
		RequirementName:    requirementName,
		PrinterID:          printerID,
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

func (s *PDFService) Reprint(ctx context.Context, printLogID int, overridePrinterID int, overridePrinterName string, userID int) (model.PDFPrintLog, error) {
	var existing model.PDFPrintLog
	existing, err := s.pdfRepo.GetPrintLog(ctx, s.db, printLogID)
	if err != nil {
		return model.PDFPrintLog{}, err
	}

	logEntry, err := s.fetchGenerationLog(ctx, existing.PDFGenerationLogID)
	if err != nil {
		return model.PDFPrintLog{}, err
	}

	printerID := existing.PrinterID
	printerName := existing.PrinterName
	if overridePrinterID != 0 {
		printerID = overridePrinterID
		printerName = overridePrinterName
	}

	return s.PrintAndLog(ctx, logEntry.TemplateName, logEntry.InputData, printerID, printerName, existing.RequirementName, userID)
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
		url, err := s.fileRepo.GetSignedDownloadURL(ctx, s.swiftConn, s.db, logs[i].FileID, 15*time.Minute)
		if err != nil {
			continue
		}
		logs[i].FileURL = url
	}

	return logs, nil
}

func (s *PDFService) FetchPDFFile(ctx context.Context, fileID string) ([]byte, string, error) {
	if fileID == "" {
		return nil, "", fmt.Errorf("file id is required")
	}
	if s.pdfRepo == nil || s.fileRepo == nil || s.swiftConn == nil || s.db == nil {
		return nil, "", fmt.Errorf("pdf service not configured for file streaming")
	}

	logEntry, err := s.pdfRepo.GetGenerationLogByFileID(ctx, s.db, fileID)
	if err != nil {
		return nil, "", err
	}

	data, err := s.fileRepo.GetFileContent(ctx, s.swiftConn, fileID)
	if err != nil {
		return nil, "", err
	}

	filename := buildPDFFilename(logEntry.PDFTitle)

	return data, filename, nil
}

func (s *PDFService) GetPDFTitleByFileID(ctx context.Context, fileID string) (string, error) {
	if fileID == "" {
		return "", fmt.Errorf("file id is required")
	}
	if s.pdfRepo == nil || s.db == nil {
		return "", fmt.Errorf("pdf service not configured for metadata lookups")
	}

	logEntry, err := s.pdfRepo.GetGenerationLogByFileID(ctx, s.db, fileID)
	if err != nil {
		return "", err
	}

	title := logEntry.PDFTitle
	if title == "" {
		title = sanitizeFilename(logEntry.TemplateName) + ".pdf"
	}
	return title, nil
}

// GeneratePDFTitle creates a downloadable/display title for a template.
func (s *PDFService) GeneratePDFTitle(templateName string) string {
	return buildPDFTitle(templateName)
}

func (s *PDFService) GeneratePDFTitleFromInput(templateName string, jsonInput []byte) string {
	return s.resolvePDFTitle(templateName, "", jsonInput)
}

// GeneratePDFFilename returns a sanitized filename (with .pdf) based on a display title.
func (s *PDFService) GeneratePDFFilename(pdfTitle string) string {
	return buildPDFFilename(pdfTitle)
}

func (s *PDFService) resolvePDFTitle(templateName, providedTitle string, jsonInput []byte) string {
	if strings.TrimSpace(providedTitle) != "" {
		return providedTitle
	}
	if tmpl, ok := pdftemplate.Registry[templateName]; ok && tmpl.TitleGenerator != nil {
		if generated, err := tmpl.TitleGenerator(jsonInput); err == nil && strings.TrimSpace(generated) != "" {
			return generated
		}
	}
	return s.GeneratePDFTitle(templateName)
}

func sanitizeFilename(name string) string {
	name = strings.ToLower(name)
	builder := strings.Builder{}
	for _, r := range name {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			builder.WriteRune(r)
		case r == ' ' || r == '_' || r == '-':
			builder.WriteRune('-')
		}
	}
	if builder.Len() == 0 {
		return "document"
	}
	return strings.Trim(builder.String(), "-")
}

func buildPDFTitle(templateName string) string {
	cleanName := strings.TrimSpace(templateName)
	if cleanName == "" {
		cleanName = "PDF"
	}
	return fmt.Sprintf("%s-%s", cleanName, time.Now().Format("200601021504"))
}

func buildPDFFilename(title string) string {
	base := sanitizeFilename(title)
	if base == "" {
		base = "document"
	}
	if strings.HasSuffix(base, ".pdf") {
		return base
	}
	return base + ".pdf"
}
