package model

import (
	"time"
)

type File struct {
	FileID      string    `json:"file_id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	SizeBytes   int       `json:"size_bytes"`
	Status      string    `json:"status"`
	Entity      string    `json:"entity"`
	EntityID    int       `json:"entity_id"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	DownloadURL string
}
