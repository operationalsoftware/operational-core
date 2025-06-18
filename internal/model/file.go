package model

import "github.com/google/uuid"

type File struct {
	FileUUID   uuid.UUID
	FileName   string
	MimeType   string
	FileExt    string
	BucketName string
}
