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
