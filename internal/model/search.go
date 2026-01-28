package model

import (
	"time"
)

type SearchInput struct {
	Q string
	E []string
}

type SearchResults struct {
	StockItems     []StockItemSearchResult
	Users          []UserSearchResult
	Resources      []ResourceSearchResult
	Services       []ServiceSearchResult
	RecentSearches []RecentSearch
}

type BaseSearchResult struct {
	Relevance int
}

type StockItemSearchResult struct {
	BaseSearchResult
	StockItemID int
	StockCode   string
	Description string
}

type UserSearchResult struct {
	BaseSearchResult
	ID        int
	Email     string
	Username  string
	FirstName string
	LastName  string
}

type ResourceSearchResult struct {
	BaseSearchResult
	ResourceID               int
	Type                     string
	Reference                string
	ServiceOwnershipTeamName string
	IsArchived               bool
}

type ServiceSearchResult struct {
	BaseSearchResult
	ResourceServiceID int
	ResourceID        int
	ResourceType      string
	ResourceReference string
	Status            ResourceServiceStatus
	StartedAt         time.Time
	StartedByUsername string
}

// Recent Search
type RecentSearch struct {
	ID             int
	SearchTerm     string
	SearchEntities []string
	UserID         int
	LastSearchedAt time.Time
}

type SearchEntity struct {
	Name          string
	Title         string
	HasPermission func(UserPermissions) bool
}
