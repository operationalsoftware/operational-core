package model

import (
	"encoding/json"
	"time"
)

type PDFGenerationLog struct {
	PDFGenerationLogID int
	TemplateName       string
	InputData          json.RawMessage
	FileID             *string
	PDFTitle           string
	PrintNodeOptions   json.RawMessage
	FileURL            *string
	CreatedByUsername  string
	CreatedAt          time.Time
}

type CreatePDFGenerationLogParams struct {
	TemplateName     string
	InputData        json.RawMessage
	PrintNodeOptions json.RawMessage
	CreatedBy        int
}

type PDFPrintLog struct {
	PDFPrintLogID      int
	PDFGenerationLogID int
	TemplateName       string
	InputData          json.RawMessage
	PrintRequirementID int
	RequirementName    string
	PrintNodeJobID     *int
	ErrorMessage       *string
	FileID             *string
	FileURL            *string
	PDFTitle           string
	CreatedByUsername  string
	CreatedAt          time.Time
}

type CreatePDFPrintLogParams struct {
	PDFGenerationLogID int
	TemplateName       string
	InputData          json.RawMessage
	PrintRequirementID int
	RequirementName    string
	CreatedBy          int
}

type PrintRequirement struct {
	PrintRequirementID int
	RequirementName    string
	PrinterName        string
	AssignedBy         int
	AssignedAt         time.Time
}
