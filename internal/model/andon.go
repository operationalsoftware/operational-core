package model

import (
	"app/pkg/appsort"
	"time"
)

type AndonStatus string

const (
	AndonStatusCancelled               AndonStatus = "Cancelled"
	AndonStatusClosed                  AndonStatus = "Closed"
	AndonStatusWorkInProgress          AndonStatus = "Work In Progress"
	AndonStatusRequiresAcknowledgement AndonStatus = "Requires Acknowledgement"
	AndonStatusOutstanding             AndonStatus = "Outstanding"
)

type Andon struct {
	AndonID             int
	Description         string
	AndonIssueID        int
	GalleryID           int
	CommentThreadID     int
	IssueName           string        `sortable:"true"`
	NamePath            []string      `sortable:"true"`
	Severity            AndonSeverity `sortable:"true"`
	Source              string        `sortable:"true"`
	Location            string        `sortable:"true"`
	IsOpen              bool          `sortable:"true"`
	Status              AndonStatus   `sortable:"true"`
	RaisedBy            int
	RaisedByUsername    string     `sortable:"true"`
	RaisedAt            time.Time  `sortable:"true"`
	ClosedAt            *time.Time `sortable:"true"`
	OpenDurationSeconds int        `sortable:"true"`
	AssignedTeam        int
	AssignedTeamName    string `sortable:"true"`

	IsAcknowledged         bool
	AcknowledgedBy         *int
	AcknowledgedByUsername *string    `sortable:"true"`
	AcknowledgedAt         *time.Time `sortable:"true"`
	IsResolved             bool
	ResolvedBy             *int
	ResolvedByUsername     *string    `sortable:"true"`
	ResolvedAt             *time.Time `sortable:"true"`
	IsCancelled            bool
	CancelledBy            *int
	CancelledByUsername    *string    `sortable:"true"`
	CancelledAt            *time.Time `sortable:"true"`
	LastUpdated            *time.Time `sortable:"true"`
	CanUserEdit            bool
	CanUserAcknowledge     bool
	CanUserResolve         bool
	CanUserCancel          bool
	CanUserReopen          bool
}

type NewAndon struct {
	Description     string
	IssueID         int
	GalleryID       int
	CommentThreadID int
	Source          string
	Location        string
	RaisedBy        string
}

type ListAndonQuery struct {
	Sort                 appsort.Sort
	DefaultSortField     string
	DefaultSortDirection appsort.Direction
	Page                 int
	PageSize             int

	StartDate                *time.Time
	EndDate                  *time.Time
	IsOpen                   *bool
	IsAcknowledged           *bool
	IsResolved               *bool
	IsCancelled              *bool
	IssueIn                  []string
	SeverityIn               []string
	StatusIn                 []string
	TeamIn                   []string
	LocationIn               []string
	RaisedByUsernameIn       []string
	AcknowledgedByUsernameIn []string
	ResolvedByUsernameIn     []string
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
	AcknowledgedBy         *int
	AcknowledgedByUsername *string
	ResolvedBy             *int
	ResolvedByUsername     *string
	CancelledBy            *int
	CancelledByUsername    *string
	ReopenedBy             *int
	ReopenedByUsername     *string
}

type AndonFilters struct {
	StartDate                *time.Time
	EndDate                  *time.Time
	IssueIn                  []string
	SeverityIn               []string
	StatusIn                 []string
	TeamIn                   []string
	LocationIn               []string
	RaisedByUsernameIn       []string
	AcknowledgedByUsernameIn []string
	ResolvedByUsernameIn     []string
}

type AndonAvailableFilters struct {
	IssueIn                  []string
	SeverityIn               []string
	StatusIn                 []string
	TeamIn                   []string
	LocationIn               []string
	RaisedByUsernameIn       []string
	AcknowledgedByUsernameIn []string
	ResolvedByUsernameIn     []string
}
