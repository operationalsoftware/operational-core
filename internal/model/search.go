package model

import (
	"time"
)

type SearchInput struct {
	Q string
	E []string
}

type SearchResults struct {
	Users          []UserSearchResult
	RecentSearches []RecentSearch
}

type BaseSearchResult struct {
	Relevance int
}

type UserSearchResult struct {
	BaseSearchResult
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
