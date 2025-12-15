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
