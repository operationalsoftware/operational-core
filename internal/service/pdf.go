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
	printNode *PrintNodeService
}

func NewPDFService(
	db *pgxpool.Pool,
	swiftConn *swift.Connection,
	fileRepo *repository.FileRepository,
	printNode *PrintNodeService,
) *PDFService {
	return &PDFService{
		db:        db,
		swiftConn: swiftConn,
		fileRepo:  fileRepo,
		printNode: printNode,
	}
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

// RecordGeneration persists the generated PDF and a log entry. It returns the log record.
func (s *PDFService) RecordGeneration(
	ctx context.Context,
	templateName string,
	inputData string,
	pdfBytes []byte,
	userID int,
) (model.PDFGenerationLog, error) {
	if s.db == nil || s.swiftConn == nil || s.fileRepo == nil {
		return model.PDFGenerationLog{}, fmt.Errorf("pdf service not fully configured for logging")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return model.PDFGenerationLog{}, err
	}
	defer tx.Rollback(ctx)

	var uid interface{}
	if userID != 0 {
		uid = userID
	}

	var logID int
	err = tx.QueryRow(ctx, `
INSERT INTO pdf_generation_log (
	template_name,
	input_data,
	created_by
) VALUES ($1, $2, $3)
RETURNING pdf_generation_log_id
`, templateName, inputData, uid).Scan(&logID)
	if err != nil {
		return model.PDFGenerationLog{}, fmt.Errorf("failed to insert pdf generation log: %w", err)
	}

	filename := fmt.Sprintf("%s-%s.pdf", sanitizeFilename(templateName), time.Now().Format("20060102-150405"))

	file, err := s.fileRepo.SaveFileContent(ctx, tx, s.swiftConn, &model.File{
		Filename:    filename,
		ContentType: "application/pdf",
		SizeBytes:   len(pdfBytes),
		Entity:      "PDFGenerationLog",
		EntityID:    logID,
		UserID:      userID,
	}, pdfBytes)
	if err != nil {
		return model.PDFGenerationLog{}, fmt.Errorf("failed to store generated pdf: %w", err)
	}

	if _, err := tx.Exec(ctx, `
UPDATE pdf_generation_log
SET file_id = $2, filename = $3
WHERE pdf_generation_log_id = $1
`, logID, file.FileID, file.Filename); err != nil {
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
		Filename:     file.Filename,
		FileURL:      downloadURL,
		UserID:       userID,
		CreatedAt:    time.Now(),
	}, nil
}

func (s *PDFService) ListRecentLogs(ctx context.Context, limit int) ([]model.PDFGenerationLog, error) {
	if s.db == nil {
		return nil, fmt.Errorf("pdf service not configured for listing logs")
	}

	rows, err := s.db.Query(ctx, `
SELECT
	pdf_generation_log_id,
	template_name,
	input_data,
	file_id,
	filename,
	created_by,
	created_at
FROM pdf_generation_log
ORDER BY created_at DESC
LIMIT $1
`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.PDFGenerationLog
	for rows.Next() {
		var log model.PDFGenerationLog
		if err := rows.Scan(
			&log.ID,
			&log.TemplateName,
			&log.InputData,
			&log.FileID,
			&log.Filename,
			&log.UserID,
			&log.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	if err := rows.Err(); err != nil {
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

func (s *PDFService) fetchGenerationLog(ctx context.Context, id int) (model.PDFGenerationLog, error) {
	var logEntry model.PDFGenerationLog
	err := s.db.QueryRow(ctx, `
SELECT
	pdf_generation_log_id,
	template_name,
	input_data,
	file_id,
	filename,
	created_by,
	created_at
FROM pdf_generation_log
WHERE pdf_generation_log_id = $1
`, id).Scan(
		&logEntry.ID,
		&logEntry.TemplateName,
		&logEntry.InputData,
		&logEntry.FileID,
		&logEntry.Filename,
		&logEntry.UserID,
		&logEntry.CreatedAt,
	)
	return logEntry, err
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

	pdfBytes, err := s.GenerateFromJSON(ctx, templateName, []byte(inputData))
	if err != nil {
		return model.PDFPrintLog{}, err
	}

	genLog, err := s.RecordGeneration(ctx, templateName, inputData, pdfBytes, userID)
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

	var printLogID int
	if s.db != nil {
		_ = s.db.QueryRow(ctx, `
INSERT INTO pdf_print_log (
	pdf_generation_log_id,
	template_name,
	requirement_name,
	printer_id,
	printer_name,
	printnode_job_id,
	status,
	error_message,
	created_by
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
RETURNING pdf_print_log_id
`, genLog.ID, templateName, requirementName, printerID, printerName, jobID, status, errorMessage, nullOrInt(userID)).Scan(&printLogID)
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
		Filename:           genLog.Filename,
		FileURL:            genLog.FileURL,
		UserID:             userID,
		CreatedAt:          time.Now(),
	}, err
}

func (s *PDFService) Reprint(ctx context.Context, printLogID int, overridePrinterID int, overridePrinterName string, userID int) (model.PDFPrintLog, error) {
	var existing model.PDFPrintLog
	err := s.db.QueryRow(ctx, `
SELECT
	pl.pdf_print_log_id,
	pl.pdf_generation_log_id,
	pl.template_name,
	pl.requirement_name,
	pl.printer_id,
	pl.printer_name
FROM pdf_print_log pl
WHERE pl.pdf_print_log_id = $1
`, printLogID).Scan(
		&existing.ID,
		&existing.PDFGenerationLogID,
		&existing.TemplateName,
		&existing.RequirementName,
		&existing.PrinterID,
		&existing.PrinterName,
	)
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
	if s.db == nil {
		return nil, fmt.Errorf("pdf service not configured for listing print logs")
	}

	rows, err := s.db.Query(ctx, `
SELECT
	pl.pdf_print_log_id,
	pl.pdf_generation_log_id,
	pl.template_name,
	pl.requirement_name,
	pl.printer_id,
	pl.printer_name,
	pl.printnode_job_id,
	pl.status,
	pl.error_message,
	pl.created_by,
	pl.created_at,
	gl.file_id,
	gl.filename,
	gl.input_data
FROM pdf_print_log pl
JOIN pdf_generation_log gl ON gl.pdf_generation_log_id = pl.pdf_generation_log_id
ORDER BY pl.created_at DESC
LIMIT $1
`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.PDFPrintLog
	for rows.Next() {
		var log model.PDFPrintLog
		if err := rows.Scan(
			&log.ID,
			&log.PDFGenerationLogID,
			&log.TemplateName,
			&log.RequirementName,
			&log.PrinterID,
			&log.PrinterName,
			&log.PrintNodeJobID,
			&log.Status,
			&log.ErrorMessage,
			&log.UserID,
			&log.CreatedAt,
			&log.FileID,
			&log.Filename,
			&log.InputData,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	if err := rows.Err(); err != nil {
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

func nullOrInt(v int) interface{} {
	if v == 0 {
		return nil
	}
	return v
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
