package model

import (
	"app/pkg/appsort"
	"time"
)

type AndonStatus string

const (
	AndonStatusCancelled      AndonStatus = "Cancelled"
	AndonStatusClosed         AndonStatus = "Closed"
	AndonStatusWorkInProgress AndonStatus = "Work In Progress"
	AndonStatusOutstanding    AndonStatus = "Outstanding"
)

type Andon struct {
	AndonID          int
	Description      string
	AndonIssueID     int
	IssueName        string        `sortable:"true"`
	NamePath         []string      `sortable:"true"`
	Severity         AndonSeverity `sortable:"true"`
	Source           string        `sortable:"true"`
	Location         string        `sortable:"true"`
	IsOpen           bool          `sortable:"true"`
	Status           AndonStatus   `sortable:"true"`
	RaisedBy         int           `sortable:"true"`
	RaisedByUsername string        `sortable:"true"`
	RaisedAt         time.Time     `sortable:"true"`
	AssignedTeam     int           `sortable:"true"`
	AssignedTeamName string        `sortable:"true"`

	IsAcknowledged         bool
	AcknowledgedBy         *int       `sortable:"true"`
	AcknowledgedByUsername *string    `sortable:"true"`
	AcknowledgedAt         *time.Time `sortable:"true"`
	IsResolved             bool
	ResolvedBy             *int       `sortable:"true"`
	ResolvedByUsername     *string    `sortable:"true"`
	ResolvedAt             *time.Time `sortable:"true"`
	IsCancelled            bool
	CancelledBy            *int       `sortable:"true"`
	CancelledByUsername    *string    `sortable:"true"`
	CancelledAt            *time.Time `sortable:"true"`
	LastUpdated            *time.Time `sortable:"true"`
	CanUserAcknowledge     bool
	CanUserResolve         bool
	CanUserCancel          bool
	CanUserReopen          bool
}

type NewAndon struct {
	Description string
	IssueID     int
	Source      string
	Location    string
	RaisedBy    string
}

type ListAndonQuery struct {
	Sort     appsort.Sort
	Page     int
	PageSize int

	StartDate              *time.Time
	EndDate                *time.Time
	IsOpen                 *bool
	IsAcknowledged         *bool
	IsResolved             *bool
	IsCancelled            *bool
	Issues                 []string
	Severities             []string
	Teams                  []string
	Locations              []string
	RaisedByUsername       []string
	AcknowledgedByUsername []string
	ResolvedByUsername     []string

	OrderBy          string
	OrderByDirection string
}

type AndonChange struct {
	AndonID                int
	AndonChangeID          int
	ChangeBy               int
	ChangeByUsername       string
	ChangeAt               time.Time
	IsCreation             bool
	Description            *string
	RaisedBy               *int
	RaisedByUsername       *string
	RaisedAt               *time.Time
	AcknowledgedBy         *int
	AcknowledgedByUsername *string
	AcknowledgedAt         *time.Time
	ResolvedBy             *int
	ResolvedByUsername     *string
	ResolvedAt             *time.Time
	CancelledBy            *int
	CancelledByUsername    *string
	CancelledAt            *time.Time
}

type AndonFilters struct {
	StartDate              *time.Time
	EndDate                *time.Time
	Issues                 []string
	Severities             []string
	Teams                  []string
	Locations              []string
	RaisedByUsername       []string
	AcknowledgedByUsername []string
	ResolvedByUsername     []string
}

type AndonAvailableFilters struct {
	IssueIn                  []string
	SeverityIn               []string
	TeamIn                   []string
	LocationIn               []string
	RaisedByUsernameIn       []string
	AcknowledgedByUsernameIn []string
	ResolvedByUsernameIn     []string
}
