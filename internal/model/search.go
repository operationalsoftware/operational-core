package model

type SearchInput struct {
	Q string
	E []string
}

// Searching entity models
type SearchResult struct {
	Type      string
	Data      interface{}
	Relevance int
}

type SearchResults struct {
	Users   *[]UserSearchResult
	Batches *[]BatchSearchResult
}

type BaseSearchResult struct {
	Relevance int `json:"relevance"`
}

type UserSearchResult struct {
	BaseSearchResult
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type BatchSearchResult struct {
	BaseSearchResult
	BatchNumber      string `json:"batch_number"`
	WorksOrderNumber string `json:"works_order_number"`
	PartNumber       string `json:"part_number"`
}
