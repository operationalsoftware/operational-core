package model

import "time"

type PDFGenerationLog struct {
	PDFGenerationLogID int
	TemplateName       string
	InputData          string
	FileID             string
	PDFTitle           string
	FileURL            string
	UserID             int
	CreatedAt          time.Time
}

type PDFPrintLog struct {
	PDFPrintLogID      int
	PDFGenerationLogID int
	TemplateName       string
	InputData          string
	RequirementName    string
	PrinterID          int
	PrinterName        string
	PrintNodeJobID     int
	Status             string
	ErrorMessage       string
	FileID             string
	PDFTitle           string
	FileURL            string
	UserID             int
	CreatedAt          time.Time
}

type PrintRequirement struct {
	PrintRequirementID int
	RequirementName    string
	PrinterID          int
	PrinterName        string
	AssignedBy         int
	AssignedAt         time.Time
}
