package repository

import (
	"app/internal/model"
	"app/pkg/db"
	"context"
)

type PDFRepository struct{}

func NewPDFRepository() *PDFRepository {
	return &PDFRepository{}
}

func (r *PDFRepository) InsertGenerationLog(
	ctx context.Context,
	exec db.PGExecutor,
	templateName string,
	inputData string,
	userID int,
	pdfTitle string,
) (int, error) {
	var id int
	err := exec.QueryRow(ctx, `
INSERT INTO pdf_generation_log (
	template_name,
	input_data,
	pdf_title,
	created_by
) VALUES ($1, $2, $3, $4)
RETURNING pdf_generation_log_id
`, templateName, inputData, pdfTitle, userID).Scan(&id)
	return id, err
}

func (r *PDFRepository) UpdateGenerationLogFile(
	ctx context.Context,
	exec db.PGExecutor,
	logID int,
	fileID string,
	pdfTitle string,
) error {
	_, err := exec.Exec(ctx, `
UPDATE pdf_generation_log
SET file_id = $2, pdf_title = $3
WHERE pdf_generation_log_id = $1
`, logID, fileID, pdfTitle)
	return err
}

func (r *PDFRepository) GetGenerationLog(
	ctx context.Context,
	exec db.PGExecutor,
	id int,
) (model.PDFGenerationLog, error) {
	var log model.PDFGenerationLog
	err := exec.QueryRow(ctx, `
SELECT
	pdf_generation_log_id,
	template_name,
	input_data,
	file_id,
	pdf_title,
	created_by,
	created_at
FROM pdf_generation_log
WHERE pdf_generation_log_id = $1
`, id).Scan(
		&log.PDFGenerationLogID,
		&log.TemplateName,
		&log.InputData,
		&log.FileID,
		&log.PDFTitle,
		&log.UserID,
		&log.CreatedAt,
	)
	return log, err
}

func (r *PDFRepository) ListRecentGenerationLogs(
	ctx context.Context,
	exec db.PGExecutor,
	limit int,
) ([]model.PDFGenerationLog, error) {
	rows, err := exec.Query(ctx, `
SELECT
	pdf_generation_log_id,
	template_name,
	input_data,
	file_id,
	pdf_title,
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
			&log.PDFGenerationLogID,
			&log.TemplateName,
			&log.InputData,
			&log.FileID,
			&log.PDFTitle,
			&log.UserID,
			&log.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, rows.Err()
}

func (r *PDFRepository) ListGenerationLogs(
	ctx context.Context,
	exec db.PGExecutor,
	limit int,
	offset int,
) ([]model.PDFGenerationLog, error) {
	rows, err := exec.Query(ctx, `
SELECT
	pdf_generation_log_id,
	template_name,
	input_data,
	file_id,
	pdf_title,
	created_by,
	created_at
FROM pdf_generation_log
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.PDFGenerationLog
	for rows.Next() {
		var log model.PDFGenerationLog
		if err := rows.Scan(
			&log.PDFGenerationLogID,
			&log.TemplateName,
			&log.InputData,
			&log.FileID,
			&log.PDFTitle,
			&log.UserID,
			&log.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}

func (r *PDFRepository) CountGenerationLogs(
	ctx context.Context,
	exec db.PGExecutor,
) (int, error) {
	var count int
	err := exec.QueryRow(ctx, `SELECT COUNT(*) FROM pdf_generation_log`).Scan(&count)
	return count, err
}

func (r *PDFRepository) InsertPrintLog(
	ctx context.Context,
	exec db.PGExecutor,
	genLogID int,
	templateName string,
	requirementName string,
	printerID int,
	printerName string,
	jobID int,
	status string,
	errorMessage string,
	userID int,
) (int, error) {
	var id int
	err := exec.QueryRow(ctx, `
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
`, genLogID, templateName, requirementName, printerID, printerName, jobID, status, errorMessage, userID).Scan(&id)
	return id, err
}

func (r *PDFRepository) GetPrintLog(
	ctx context.Context,
	exec db.PGExecutor,
	printLogID int,
) (model.PDFPrintLog, error) {
	var log model.PDFPrintLog
	err := exec.QueryRow(ctx, `
SELECT
	pdf_print_log_id,
	pdf_generation_log_id,
	template_name,
	requirement_name,
	printer_id,
	printer_name,
	printnode_job_id,
	status,
	error_message,
	created_by,
	created_at
FROM pdf_print_log
WHERE pdf_print_log_id = $1
`, printLogID).Scan(
		&log.PDFPrintLogID,
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
	)
	return log, err
}

func (r *PDFRepository) ListRecentPrintLogs(
	ctx context.Context,
	exec db.PGExecutor,
	limit int,
) ([]model.PDFPrintLog, error) {
	rows, err := exec.Query(ctx, `
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
	gl.pdf_title,
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
			&log.PDFPrintLogID,
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
			&log.PDFTitle,
			&log.InputData,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}

func (r *PDFRepository) ListPrintRequirements(
	ctx context.Context,
	exec db.PGExecutor,
) ([]model.PrintRequirement, error) {
	rows, err := exec.Query(ctx, `
SELECT
	print_requirement_id,
	requirement_name,
	COALESCE(printer_id, 0) AS printer_id,
	COALESCE(printer_name, '') AS printer_name,
	COALESCE(assigned_by, 0) AS assigned_by,
	COALESCE(assigned_at, NOW()) AS assigned_at
FROM pdf_print_requirement
ORDER BY requirement_name ASC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reqs []model.PrintRequirement
	for rows.Next() {
		var pr model.PrintRequirement
		if err := rows.Scan(
			&pr.PrintRequirementID,
			&pr.RequirementName,
			&pr.PrinterID,
			&pr.PrinterName,
			&pr.AssignedBy,
			&pr.AssignedAt,
		); err != nil {
			return nil, err
		}
		reqs = append(reqs, pr)
	}
	return reqs, rows.Err()
}

func (r *PDFRepository) UpsertPrintRequirement(
	ctx context.Context,
	exec db.PGExecutor,
	req model.PrintRequirement,
) (model.PrintRequirement, error) {
	err := exec.QueryRow(ctx, `
INSERT INTO pdf_print_requirement (
	requirement_name,
	printer_id,
	printer_name,
	assigned_by
) VALUES ($1,$2,$3,$4)
ON CONFLICT (requirement_name)
DO UPDATE SET
	printer_id = EXCLUDED.printer_id,
	printer_name = EXCLUDED.printer_name,
	assigned_by = EXCLUDED.assigned_by,
	assigned_at = NOW()
RETURNING
	print_requirement_id,
	assigned_at
`, req.RequirementName, req.PrinterID, req.PrinterName, req.AssignedBy).Scan(
		&req.PrintRequirementID,
		&req.AssignedAt,
	)
	return req, err
}
