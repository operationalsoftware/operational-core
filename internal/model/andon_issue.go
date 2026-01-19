package model

import (
	"app/pkg/appsort"
	"time"
)

const MaxAndonIssueDepth = 4

type AndonSeverity string

const (
	AndonSeverityInfo                 AndonSeverity = "Info"
	AndonSeveritySelfResolvable       AndonSeverity = "Self-resolvable"
	AndonSeverityRequiresIntervention AndonSeverity = "Requires Intervention"
)

var AndonSeverities = []AndonSeverity{
	AndonSeverityInfo,
	AndonSeveritySelfResolvable,
	AndonSeverityRequiresIntervention,
}

type AndonIssue struct {
	AndonIssueID           int
	IssueName              string   `sortable:"true"`
	NamePath               []string `sortable:"true"`
	IsArchived             bool     `sortable:"true"`
	ChildrenCount          int      `sortable:"true"`
	Depth                  int
	ParentID               *int
	AssignedTeam           int
	AssignedTeamName       string        `sortable:"true"`
	Severity               AndonSeverity `sortable:"true"`
	RequireAcknowledgement bool

	CreatedAt         time.Time `sortable:"true"`
	CreatedBy         int
	CreatedByUsername string     `sortable:"true"`
	UpdatedAt         *time.Time `sortable:"true"`
	UpdatedBy         *int
	UpdatedByUsername *string
}

type NewAndonIssue struct {
	IssueName              string
	ParentID               *int
	AssignedTeam           int
	Severity               AndonSeverity
	RequireAcknowledgement bool
}

type AndonIssueUpdate struct {
	IssueName              string
	ParentID               *int
	IsArchived             bool
	AssignedTeam           int
	Severity               AndonSeverity
	RequireAcknowledgement bool
}

type AndonIssueGroupUpdate struct {
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

type AndonIssueNode struct {
	AndonIssueID           int
	IssueName              string   `sortable:"true"`
	NamePath               []string `sortable:"true"`
	IsArchived             bool     `sortable:"true"`
	ChildrenCount          int      `sortable:"true"`
	Depth                  int
	ParentID               *int
	AssignedTeam           *int
	AssignedTeamName       *string        `sortable:"true"`
	Severity               *AndonSeverity `sortable:"true"`
	DownDepth              int
	RequireAcknowledgement bool

	IsGroup bool

	CreatedAt         time.Time
	CreatedBy         int
	CreatedByUsername string
	UpdatedAt         *time.Time
	UpdatedBy         *int
	UpdatedByUsername *string
}

type AndonIssueGroup struct {
	AndonIssueID int
	IsArchived   bool
	IssueName    string
	ParentID     *int
	NamePath     []string
	Depth        int
	DownDepth    int
}

type NewAndonIssueGroup struct {
	IssueName string
	ParentID  *int
}
