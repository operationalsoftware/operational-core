package model

import (
	"time"
)

type SearchInput struct {
	Q string
	E []string
}

// Searching entity models
// type SearchResult struct {
// 	Type      string
// 	Data      interface{}
// 	Relevance int
// }

type SearchResults struct {
	Users          []UserSearchResult
	Batches        []BatchSearchResult
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

type BatchSearchResult struct {
	BaseSearchResult
	BatchNumber      string
	WorksOrderNumber string
	PartNumber       string
}

// Recent Search
type RecentSearch struct {
	ID             int
	SearchTerm     string
	SearchEntities []string
	UserID         int
	LastSearchedAt time.Time
}
