package model

import (
	"app/pkg/appsort"
	"time"
)

type AndonIssue struct {
	AndonIssueID  int
	IssueName     string   `sortable:"true"`
	NamePath      []string `sortable:"true"`
	IsArchived    bool     `sortable:"true"`
	ChildrenCount int      `sortable:"true"`
	Depth         int
	ParentID      *int

	CreatedAt time.Time
	CreatedBy int
	UpdatedAt *time.Time
	UpdatedBy *int
}

type NewAndonIssue struct {
	IssueName string
	ParentID  *int
}

type AndonIssueUpdate struct {
	IssueName  string
	ParentID   *int
	IsArchived bool
}

type ListAndonIssuesQuery struct {
	ShowArchived bool
	Sort         appsort.Sort
	Page         int
	PageSize     int
}
