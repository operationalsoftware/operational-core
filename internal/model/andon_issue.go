package model

import (
	"app/pkg/appsort"
	"time"
)

type AndonSeverity string

const (
	AndonSeverityInfo                 AndonSeverity = "Info"
	AndonSeveritySelfResolvable       AndonSeverity = "Self-resolvable"
	AndonSeverityRequiresIntervention AndonSeverity = "Requires Intervention"
)

type AndonIssue struct {
	AndonIssueID       int
	IssueName          string   `sortable:"true"`
	NamePath           []string `sortable:"true"`
	IsArchived         bool     `sortable:"true"`
	ChildrenCount      int      `sortable:"true"`
	Depth              int
	ParentID           *int
	AssignedToTeam     int
	AssignedToTeamName string        `sortable:"true"`
	Severity           AndonSeverity `sortable:"true"`

	CreatedAt         time.Time
	CreatedBy         int
	CreatedByUsername string
	UpdatedAt         *time.Time
	UpdatedBy         *int
	UpdatedByUsername *string
}

type NewAndonIssue struct {
	IssueName          string
	ParentID           *int
	AssignedToTeam     int
	ResolvableByRaiser bool
	WillStopProcess    bool
}

type AndonIssueUpdate struct {
	IssueName          string
	ParentID           *int
	IsArchived         bool
	AssignedToTeam     int
	ResolvableByRaiser bool
	WillStopProcess    bool
}

type ListAndonIssuesQuery struct {
	ShowArchived bool
	Sort         appsort.Sort
	Page         int
	PageSize     int
}
