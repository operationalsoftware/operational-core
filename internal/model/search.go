package model

// Searching entity models
type SearchResult struct {
	Type      string      `json:"type"` // "User", "Stock", etc.
	Data      interface{} `json:"data"`
	Relevance int         `json:"relevance"`
}

// type SearchResult struct {
// 	Type  string             `json:"type"` // "User", "Stock", etc.
// 	User  *UserSearchResult  `json:"user,omitempty"`
// 	Batch *BatchSearchResult `json:"batch,omitempty"`
// 	// more types as needed
// 	Relevance int `json:"relevance"`
// }

type UserSearchResult struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Relevance int    `json:"relevance"`
}

type BatchSearchResult struct {
	BatchNumber      string `json:"batch_number"`
	WorksOrderNumber string `json:"works_order_number"`
	PartNumber       string `json:"part_number"`
	Relevance        int    `json:"relevance"`
}
