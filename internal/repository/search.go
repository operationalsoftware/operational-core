package repository

import (
	"app/internal/model"
	"app/pkg/db"
	"context"
	"sort"
)

type SearchRepository struct{}

func NewSearchRepository() *SearchRepository {
	return &SearchRepository{}
}

func (r *SearchRepository) CreateRecentSearch(
	ctx context.Context,
	exec db.PGExecutor,
	searchInput model.SearchInput,
	userID int,
) error {
	sort.Strings(searchInput.E)

	query := `
INSERT INTO recent_search (
	search_term,
	search_entities,
	user_id
)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, search_term, search_entities)
DO UPDATE SET last_searched_at = NOW()
	`
	_, err := exec.Exec(
		ctx,
		query,
		searchInput.Q,
		searchInput.E,
		userID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *SearchRepository) FetchRecentSearches(
	ctx context.Context,
	exec db.PGExecutor,
	searchInput model.SearchInput,
	userID int,
) ([]model.RecentSearch, error) {
	var searches []model.RecentSearch

	query := `
	SELECT 
		recent_search_id,
		search_term, 
		search_entities, 
		user_id, 
		last_searched_at
	FROM 
		recent_search
	WHERE 
		user_id = $1
	ORDER BY 
		last_searched_at DESC
		LIMIT 10
	`

	rows, err := exec.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var search model.RecentSearch
		if err := rows.Scan(
			&search.ID,
			&search.SearchTerm,
			&search.SearchEntities,
			&search.UserID,
			&search.LastSearchedAt,
		); err != nil {
			return nil, err
		}
		searches = append(searches, search)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return searches, nil
}
