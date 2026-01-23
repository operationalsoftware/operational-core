package service

import (
	"app/internal/model"
	"app/internal/pdftemplate"
	"app/internal/repository"
	"app/pkg/pdf"
	"app/pkg/printnode"
	"context"
	"encoding/json"
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
	title := strings.TrimSpace(pdfTitle)
	if title == "" {
		title = pdftemplate.FallbackTitle(templateName)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return model.PDFGenerationLog{}, err
	}
	defer tx.Rollback(ctx)

	printNodeOptions := json.RawMessage(`{}`)
	logParams := model.CreatePDFGenerationLogParams{
		TemplateName:     templateName,
		InputData:        json.RawMessage(inputData),
		PrintNodeOptions: printNodeOptions,
		CreatedBy:        userID,
	}
	logID, err := s.pdfRepo.InsertGenerationLog(ctx, tx, logParams, title)
	if err != nil {
		return model.PDFGenerationLog{}, fmt.Errorf("failed to insert pdf generation log: %w", err)
	}

	fileName := s.GeneratePDFFilename(title)

	fileInput := &model.File{
		Filename:    fileName,
		ContentType: "application/pdf",
		SizeBytes:   len(pdfBytes),
		Entity:      "PDFGenerationLog",
		EntityID:    logID,
		UserID:      userID,
	}
	file, err := s.fileRepo.SaveFileContent(ctx, tx, s.swiftConn, fileInput, pdfBytes)
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

	fileID := file.FileID
	return model.PDFGenerationLog{
		PDFGenerationLogID: logID,
		TemplateName:       templateName,
		InputData:          json.RawMessage(inputData),
		FileID:             &fileID,
		PDFTitle:           title,
		PrintNodeOptions:   printNodeOptions,
		FileURL:            &downloadURL,
		CreatedAt:          time.Now(),
	}, nil
}

func (s *PDFService) ListGenerationLogs(ctx context.Context, limit, offset int) ([]model.PDFGenerationLog, int, error) {
	total, err := s.pdfRepo.CountGenerationLogs(ctx, s.db)
	if err != nil {
		return nil, 0, err
	}

	logs, err := s.pdfRepo.ListGenerationLogs(ctx, s.db, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	for i := range logs {
		if logs[i].FileID == nil || *logs[i].FileID == "" {
			continue
		}
		if logs[i].PDFTitle == "" {
			continue
		}
		url, err := s.fileRepo.GetSignedDownloadURLWithDisposition(
			ctx,
			s.swiftConn,
			s.db,
			*logs[i].FileID,
			15*time.Minute,
			s.GeneratePDFFilename(logs[i].PDFTitle),
			true,
		)
		if err != nil {
			continue
		}
		logs[i].FileURL = &url
	}

	return logs, total, nil
}

func (s *PDFService) PrintAndLog(
	ctx context.Context,
	templateName string,
	inputData string,
	requirementName string,
	userID int,
) (int, error) {
	genLog, pdfBytes, err := s.GenerateAndLog(ctx, templateName, inputData, userID)
	if err != nil {
		return 0, err
	}
	return s.PrintGeneratedPDF(ctx, genLog, pdfBytes, requirementName, "", userID)
}

func (s *PDFService) GenerateAndLog(
	ctx context.Context,
	templateName string,
	inputData string,
	userID int,
) (model.PDFGenerationLog, []byte, error) {
	pdfBytes, title, err := s.GenerateFromJSON(ctx, templateName, []byte(inputData))
	if err != nil {
		return model.PDFGenerationLog{}, nil, err
	}

	genLog, err := s.RecordGeneration(ctx, templateName, inputData, pdfBytes, userID, title)
	if err != nil {
		return model.PDFGenerationLog{}, nil, err
	}

	return genLog, pdfBytes, nil
}

func (s *PDFService) PrintGeneratedPDF(
	ctx context.Context,
	genLog model.PDFGenerationLog,
	pdfBytes []byte,
	requirementName string,
	overridePrinterName string,
	userID int,
) (int, error) {
	requirementName = strings.TrimSpace(requirementName)
	if requirementName == "" {
		return 0, fmt.Errorf("print requirement name is required")
	}

	requirement, err := s.pdfRepo.GetPrintRequirementByName(ctx, s.db, requirementName)
	if err != nil {
		return 0, err
	}

	printerName := strings.TrimSpace(overridePrinterName)
	if printerName == "" {
		printerName = strings.TrimSpace(requirement.PrinterName)
	}
	if printerName == "" {
		return 0, fmt.Errorf("print requirement does not have a printer assigned: %s", requirementName)
	}

	printLogParams := model.CreatePDFPrintLogParams{
		PDFGenerationLogID: genLog.PDFGenerationLogID,
		TemplateName:       genLog.TemplateName,
		InputData:          genLog.InputData,
		PrintRequirementID: requirement.PrintRequirementID,
		RequirementName:    requirement.RequirementName,
		CreatedBy:          userID,
	}

	printerID := 0
	printers, err := s.printNode.Printers(ctx)
	if err != nil {
		return s.recordPrintError(ctx, printLogParams, err)
	}

	for _, pr := range printers {
		if strings.EqualFold(pr.Name, printerName) {
			printerID = pr.ID
			break
		}
	}
	if printerID == 0 {
		err = fmt.Errorf("printer not found: %s", printerName)
		return s.recordPrintError(ctx, printLogParams, err)
	}

	jobID, err := s.printNode.SubmitPDF(ctx, printerID, genLog.TemplateName, pdfBytes)
	if err != nil {
		return s.recordPrintError(ctx, printLogParams, err)
	}

	printLogID := 0
	var jobIDPtr *int
	if jobID != 0 {
		jobIDPtr = &jobID
	}
	id, insertErr := s.pdfRepo.InsertPrintLog(ctx, s.db, printLogParams, jobIDPtr, nil)
	if insertErr == nil {
		printLogID = id
	}
	return printLogID, err
}

func (s *PDFService) Reprint(ctx context.Context, printLogID int, overridePrinterName string, userID int) (int, error) {
	var existing model.PDFPrintLog
	existing, err := s.pdfRepo.GetPrintLog(ctx, s.db, printLogID)
	if err != nil {
		return 0, err
	}

	logEntry, err := s.pdfRepo.GetGenerationLog(ctx, s.db, existing.PDFGenerationLogID)
	if err != nil {
		return 0, err
	}

	var requirementName string
	req, err := s.pdfRepo.GetPrintRequirementByID(ctx, s.db, existing.PrintRequirementID)
	if err != nil {
		return 0, err
	}
	requirementName = req.RequirementName
	if requirementName == "" {
		return 0, fmt.Errorf("print requirement not available for reprint")
	}

	genLog, pdfBytes, err := s.GenerateAndLog(ctx, logEntry.TemplateName, string(logEntry.InputData), userID)
	if err != nil {
		return 0, err
	}

	overridePrinterName = strings.TrimSpace(overridePrinterName)
	return s.PrintGeneratedPDF(ctx, genLog, pdfBytes, requirementName, overridePrinterName, userID)
}

func (s *PDFService) ListRecentPrintLogs(ctx context.Context, limit int) ([]model.PDFPrintLog, error) {
	logs, err := s.pdfRepo.ListRecentPrintLogs(ctx, s.db, limit)
	if err != nil {
		return nil, err
	}

	for i := range logs {
		if logs[i].FileID == nil || *logs[i].FileID == "" {
			continue
		}
		if logs[i].PDFTitle == "" {
			continue
		}
		url, err := s.fileRepo.GetSignedDownloadURLWithDisposition(
			ctx,
			s.swiftConn,
			s.db,
			*logs[i].FileID,
			15*time.Minute,
			s.GeneratePDFFilename(logs[i].PDFTitle),
			true,
		)
		if err != nil {
			continue
		}
		logs[i].FileURL = &url
	}

	return logs, nil
}

func (s *PDFService) ListPrintRequirements(ctx context.Context) ([]model.PrintRequirement, error) {
	return s.pdfRepo.ListPrintRequirements(ctx, s.db)
}

func (s *PDFService) GetPrintRequirementByName(ctx context.Context, requirementName string) (model.PrintRequirement, error) {
	if requirementName == "" {
		return model.PrintRequirement{}, fmt.Errorf("print requirement name is required")
	}
	return s.pdfRepo.GetPrintRequirementByName(ctx, s.db, requirementName)
}

func (s *PDFService) SavePrintRequirement(ctx context.Context, requirementName string, printerName string, assignedBy int) (model.PrintRequirement, error) {
	requirementName = strings.TrimSpace(requirementName)
	if requirementName == "" {
		return model.PrintRequirement{}, fmt.Errorf("requirement name is required")
	}
	printerName = strings.TrimSpace(printerName)
	if printerName == "" {
		return model.PrintRequirement{}, fmt.Errorf("printer name is required")
	}

	assignments, err := s.pdfRepo.ListPrintRequirements(ctx, s.db)
	if err != nil {
		return model.PrintRequirement{}, err
	}

	for _, a := range assignments {
		if printerName == "" {
			break
		}
		if strings.EqualFold(a.PrinterName, printerName) && !strings.EqualFold(a.RequirementName, requirementName) {
			return model.PrintRequirement{}, fmt.Errorf("printer already assigned to another requirement")
		}
	}

	return s.pdfRepo.UpdatePrintRequirement(ctx, s.db, requirementName, printerName, assignedBy)
}

func (s *PDFService) ListAvailablePrinters(ctx context.Context, currentReq string, printers []printnode.Printer) ([]printnode.Printer, error) {
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

// helper to record a print error and return the log id
func (s *PDFService) recordPrintError(
	ctx context.Context,
	params model.CreatePDFPrintLogParams,
	err error,
) (int, error) {
	errorMessage := err.Error()
	errorPtr := &errorMessage
	logID, _ := s.pdfRepo.InsertPrintLog(ctx, s.db, params, nil, errorPtr)
	return logID, err
}

// GeneratePDFFilename returns a filename (with .pdf) based on the PDF title.
func (s *PDFService) GeneratePDFFilename(pdfTitle string) string {
	name := strings.TrimSpace(pdfTitle)
	if !strings.HasSuffix(strings.ToLower(name), ".pdf") {
		name += ".pdf"
	}
	return name
}
