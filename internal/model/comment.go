package model

import (
	"time"
)

type Comment struct {
	CommentID           int
	EntityID            int
	Entity              string
	Comment             string
	CommentedBy         string
	CommentedByUsername string
	CommentedAt         time.Time
	Attachments         []File
}

type NewComment struct {
	Entity   string
	EntityID int
	Comment  string
}
