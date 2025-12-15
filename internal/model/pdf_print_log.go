package model

import "time"

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
