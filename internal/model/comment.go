package model

import (
	"time"
)

type Comment struct {
	CommentID           int
	CommentThreadID     int
	Comment             string
	CommentedBy         string
	CommentedByUsername string
	CommentedAt         time.Time
	Attachments         []File
}

type NewComment struct {
	CommentThreadID int
	Comment         string
}
