package model

import (
	"time"
)

type File struct {
	FileID           string
	ObjectName       string
	OriginalFilename string
	ContentType      string
	SizeBytes        int64
	Entity           string
	UserID           int
	CreatedAt        time.Time
}
