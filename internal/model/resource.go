package model

import (
	"app/pkg/appsort"
	"time"

	"github.com/shopspring/decimal"
)

type Resource struct {
	ResourceID               int
	Type                     string `sortable:"true"`
	Reference                string `sortable:"true"`
	ServiceOwnershipTeamID   *int
	ServiceOwnershipTeamName *string
	ServiceScheduleNames     []string
	IsArchived               bool
	LastServicedAt           *time.Time `sortable:"true"`
}

type NewResource struct {
	Type                   string
	Reference              string
	ServiceOwnershipTeamID *int
}

type ResourceUpdate struct {
	Type                   string
	Reference              string
	IsArchived             bool
	ServiceOwnershipTeamID *int
}

type NewResourceServiceMetricRecord struct {
	ResourceID              int
	ResourceServiceMetricID int
	Value                   decimal.Decimal
	ClosedByServiceID       *int
}

type ResourceFilters struct {
	TypeIn                 []string
	ServiceOwnershipTeamIn []string
	ReferenceIn            []string
}

type ResourceAvailableFilters struct {
	TypeIn                 []string
	ServiceOwnershipTeamIn []string
	ReferenceIn            []string
}

type GetResourcesQuery struct {
	Sort                   appsort.Sort
	Page                   int
	PageSize               int
	IsArchived             bool
	TypeIn                 []string
	ServiceOwnershipTeamIn []string
	ReferenceIn            []string
}

type GetServicesQuery struct {
	Sort     appsort.Sort
	Page     int
	PageSize int

	ResourceIn []string
}
