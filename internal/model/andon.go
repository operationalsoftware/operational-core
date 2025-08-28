package model

import (
	"app/pkg/appsort"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type AndonEvent struct {
	AndonEventID           int
	IssueDescription       string
	IssueID                int
	IssueName              string `sortable:"true"`
	NamePath               []string
	CanUserAcknowledge     bool
	CanUserResolve         bool
	CanUserCancel          bool
	Severity               AndonSeverity `sortable:"true"`
	Source                 string        `sortable:"true"`
	Location               string        `sortable:"true"`
	RaisedBy               int           `sortable:"true"`
	RaisedByUsername       string        `sortable:"true"`
	RaisedAt               time.Time     `sortable:"true"`
	AssignedTeam           int           `sortable:"true"`
	AssignedTeamName       string        `sortable:"true"`
	AcknowledgedBy         *int          `sortable:"true"`
	AcknowledgedByUsername *string       `sortable:"true"`
	AcknowledgedAt         *time.Time    `sortable:"true"`
	ResolvedBy             *int          `sortable:"true"`
	ResolvedByUsername     *string       `sortable:"true"`
	ResolvedAt             *time.Time    `sortable:"true"`
	CancelledBy            *int          `sortable:"true"`
	CancelledByUsername    *string       `sortable:"true"`
	CancelledAt            *time.Time    `sortable:"true"`
	Status                 string        `sortable:"true"`
	LastUpdated            *time.Time    `sortable:"true"`
}

type NewAndonEvent struct {
	IssueDescription string
	IssueID          int
	Source           string
	Location         string
	RaisedBy         string
	LinkedEntityID   int
	LinkedEntityType string
}

type AndonEventUpdate struct {
	IssueDescription string
	IssueID          int
	Source           string
	Location         string
	LastUpdated      time.Time
}

type ListAndonQuery struct {
	Sort     appsort.Sort
	Page     int
	PageSize int

	StartDate              *time.Time
	EndDate                *time.Time
	Issues                 []string
	Serverities            []string
	Teams                  []string
	Locations              []string
	Statuses               []string
	RaisedByUsername       []string
	AcknowledgedByUsername []string
	ResolvedByUsername     []string

	OrderBy          string
	OrderByDirection string
}

type AndonChange struct {
	AndonEventID           int
	IssueDescription       pgtype.Text
	IssueID                pgtype.Int4
	Location               pgtype.Text
	RaisedBy               pgtype.Text
	RaisedByUsername       pgtype.Text
	RaisedAt               pgtype.Timestamptz
	AcknowledgedBy         pgtype.Int4
	AcknowledgedByUsername pgtype.Text
	AcknowledgedAt         pgtype.Timestamptz
	ResolvedBy             pgtype.Int4
	ResolvedByUsername     pgtype.Text
	ResolvedAt             pgtype.Timestamptz
	CancelledBy            pgtype.Int4
	CancelledByUsername    pgtype.Text
	CancelledAt            pgtype.Timestamptz
	ChangeBy               int
	ChangeByUsername       string
	ChangeAt               time.Time
	IsCreation             bool
	Status                 string
	LastUpdated            pgtype.Timestamptz
}

type AndonFilters struct {
	StartDate              *time.Time
	EndDate                *time.Time
	Issues                 []string
	Severities             []string
	Teams                  []string
	Locations              []string
	Statuses               []string
	RaisedByUsername       []string
	AcknowledgedByUsername []string
	ResolvedByUsername     []string
}

type AndonAvailableFilters struct {
	IssueIn                  []string
	SeverityIn               []string
	TeamIn                   []string
	LocationIn               []string
	StatusIn                 []string
	RaisedByUsernameIn       []string
	AcknowledgedByUsernameIn []string
	ResolvedByUsernameIn     []string
}
