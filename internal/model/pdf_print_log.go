package model

import "time"

type PDFPrintLog struct {
	ID                 int
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
	Filename           string
	FileURL            string
	UserID             int
	CreatedAt          time.Time
}
