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
	input model.CreatePDFGenerationLogParams,
	pdfTitle string,
) (int, error) {
	var id int
	err := exec.QueryRow(ctx, `
INSERT INTO pdf_generation_log (
	template_name,
	input_data,
	pdf_title,
	print_node_options,
	created_by
) VALUES ($1, $2::jsonb, $3, $4::jsonb, $5)
RETURNING pdf_generation_log_id
`, input.TemplateName, input.InputData, pdfTitle, input.PrintNodeOptions, input.CreatedBy).Scan(&id)
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
	print_node_options,
	created_by_username,
	created_at
FROM pdf_generation_log_view
WHERE pdf_generation_log_id = $1
`, id).Scan(
		&log.PDFGenerationLogID,
		&log.TemplateName,
		&log.InputData,
		&log.FileID,
		&log.PDFTitle,
		&log.PrintNodeOptions,
		&log.CreatedByUsername,
		&log.CreatedAt,
	)
	return log, err
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
	print_node_options,
	created_by_username,
	created_at
FROM pdf_generation_log_view
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
			&log.PrintNodeOptions,
			&log.CreatedByUsername,
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
	input model.CreatePDFPrintLogParams,
	jobID *int64,
	errorMessage *string,
) (int, error) {
	var id int
	err := exec.QueryRow(ctx, `
INSERT INTO pdf_print_log (
	pdf_generation_log_id,
	template_name,
	print_requirement_id,
	printnode_job_id,
	error_message,
	created_by
) VALUES ($1,$2,$3,$4,$5,$6)
RETURNING pdf_print_log_id
`, input.PDFGenerationLogID, input.TemplateName, input.PrintRequirementID, jobID, errorMessage, input.CreatedBy).Scan(&id)
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
	print_requirement_id,
	printnode_job_id,
	error_message,
	created_at
FROM pdf_print_log
WHERE pdf_print_log_id = $1
`, printLogID).Scan(
		&log.PDFPrintLogID,
		&log.PDFGenerationLogID,
		&log.TemplateName,
		&log.PrintRequirementID,
		&log.PrintNodeJobID,
		&log.ErrorMessage,
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
	pdf_print_log_id,
	pdf_generation_log_id,
	template_name,
	requirement_name,
	printnode_job_id,
	error_message,
	created_by_username,
	created_at,
	file_id,
	pdf_title,
	input_data
FROM pdf_print_log_view
ORDER BY created_at DESC
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
			&log.PrintNodeJobID,
			&log.ErrorMessage,
			&log.CreatedByUsername,
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
	printer_name,
	assigned_by,
	assigned_at
FROM print_requirement
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

func (r *PDFRepository) GetPrintRequirementByName(
	ctx context.Context,
	exec db.PGExecutor,
	requirementName string,
) (model.PrintRequirement, error) {
	var pr model.PrintRequirement
	err := exec.QueryRow(ctx, `
SELECT
	print_requirement_id,
	requirement_name,
	printer_name,
	assigned_by,
	assigned_at
FROM print_requirement
WHERE requirement_name = $1
`, requirementName).Scan(
		&pr.PrintRequirementID,
		&pr.RequirementName,
		&pr.PrinterName,
		&pr.AssignedBy,
		&pr.AssignedAt,
	)
	return pr, err
}

func (r *PDFRepository) GetPrintRequirementByID(
	ctx context.Context,
	exec db.PGExecutor,
	id int,
) (model.PrintRequirement, error) {
	var pr model.PrintRequirement
	err := exec.QueryRow(ctx, `
SELECT
	print_requirement_id,
	requirement_name,
	printer_name,
	assigned_by,
	assigned_at
FROM print_requirement
WHERE print_requirement_id = $1
`, id).Scan(
		&pr.PrintRequirementID,
		&pr.RequirementName,
		&pr.PrinterName,
		&pr.AssignedBy,
		&pr.AssignedAt,
	)
	return pr, err
}

func (r *PDFRepository) UpdatePrintRequirement(
	ctx context.Context,
	exec db.PGExecutor,
	requirementName string,
	printerName string,
	assignedBy int,
) (model.PrintRequirement, error) {
	req := model.PrintRequirement{
		RequirementName: requirementName,
		PrinterName:     printerName,
		AssignedBy:      assignedBy,
	}
	err := exec.QueryRow(ctx, `
UPDATE print_requirement
SET
	printer_name = $2,
	assigned_by = $3,
	assigned_at = NOW()
WHERE requirement_name = $1
RETURNING
	print_requirement_id,
	assigned_at
`, req.RequirementName, req.PrinterName, req.AssignedBy).Scan(
		&req.PrintRequirementID,
		&req.AssignedAt,
	)
	return req, err
}
