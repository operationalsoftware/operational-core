package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type ResourceServiceStatus string

const (
	ServiceStatusCancelled      ResourceServiceStatus = "Cancelled"
	ServiceStatusCompleted      ResourceServiceStatus = "Completed"
	ServiceStatusWorkInProgress ResourceServiceStatus = "Work In Progress"
)

type ServiceMetric struct {
	ServiceMetricID int
	Name            string `sortable:"true"`
	Description     string
	IsCumulative    bool `sortable:"true"`
	IsArchived      bool `sortable:"true"`
}

type ServiceSchedule struct {
	ServiceScheduleID       int
	Name                    string `sortable:"true"`
	ResourceServiceMetricID int
	MetricName              string          `sortable:"true"`
	Threshold               decimal.Decimal `sortable:"true"`
	IsArchived              bool            `sortable:"true"`
}

type ResourceService struct {
	ResourceServiceID   int
	ResourceID          int
	ResourceType        string
	ResourceReference   string
	Status              ResourceServiceStatus
	StartedBy           int
	StartedAt           time.Time
	CompletedBy         *int
	CompletedAt         *time.Time
	CancelledBy         *int
	CancelledAt         *time.Time
	StartedByUsername   string
	CompletedByUsername *string
	CancelledByUsername *string
	Notes               string
	GalleryID           int
	GalleryURL          string
	CommentThreadID     int
}

type ResourceServiceChange struct {
	ResourceServiceID       int
	ResourceServiceChangeID int
	ChangeBy                int
	ChangeByUsername        string
	ChangeAt                time.Time
	IsCreation              bool
	Notes                   *string
	StartedBy               *int
	StartedByUsername       *string
	CompletedBy             *int
	CompletedByUsername     *string
	ReopenedBy              *int
	ReopenedByUsername      *string
	CancelledBy             *int
	CancelledByUsername     *string
}

type NewResourceService struct {
	ResourceID      int
	GalleryID       int
	CommentThreadID int
	Notes           string
}

type UpdateResourceService struct {
	ResourceServiceID int
	Notes             string
}

type ResourceServiceSchedule struct {
	ResourceServiceScheduleID int
	ResourceID                int
	ServiceScheduleID         int
	MetricName                string
	Threshold                 decimal.Decimal
	CreatedAt                 time.Time
}

type NewResourceServiceSchedule struct {
	ResourceID        int
	ServiceScheduleID int
}

type NewServiceSchedule struct {
	Name                    string
	ResourceServiceMetricID int
	Threshold               decimal.Decimal
}

type UpdateServiceSchedule struct {
	ServiceScheduleID       int
	Name                    string
	ResourceServiceMetricID int
	Threshold               decimal.Decimal
	IsArchived              bool
}

type ResourceServiceMetricStatus struct {
	ServiceScheduleID        int
	ServiceScheduleName      string
	ResourceID               int
	Type                     string
	Reference                string
	ServiceOwnershipTeamID   *int
	ServiceOwnershipTeamName *string
	ResourceServiceMetricID  int
	MetricName               string
	CurrentValue             decimal.Decimal
	Threshold                decimal.Decimal
	NormalisedValue          decimal.Decimal
	NormalisedPercentage     decimal.Decimal
	IsDue                    bool
	WIPServiceID             *int
	HasWIPService            bool
	LastRecordedAt           *time.Time
	LastServicedAt           *time.Time
	ScheduleIsArchived       bool
	MetricIsArchived         bool
}

type ResourceServiceMetricStatusesQuery struct {
	ServiceOwnershipTeamIDs []int
	Page                    int
	PageSize                int
}

type NewServiceMetric struct {
	Name         string
	Description  string
	IsCumulative bool
}

type UpdateResourceServiceMetric struct {
	ServiceMetricID int
	Name            string
	Description     string
	IsCumulative    bool
	IsArchived      bool
}
