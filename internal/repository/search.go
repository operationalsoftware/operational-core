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

func (r *SearchRepository) SearchUsers(
	ctx context.Context,
	exec db.PGExecutor,
	searchTerm string,
) ([]model.UserSearchResult, error) {

	rows, err := exec.Query(ctx, `
		SELECT
			user_id,
			COALESCE(email, '') AS email,
			COALESCE(username, '') AS username,
			COALESCE(first_name, '') AS first_name,
			COALESCE(last_name, '') AS last_name,
			(
				CASE WHEN COALESCE(email, '') ILIKE $1 THEN 3
					WHEN COALESCE(email, '') ILIKE $1 || '%' THEN 2
					WHEN COALESCE(email, '') ILIKE '%' || $1 || '%' THEN 1
					ELSE 0
				END * 3
				+
				CASE WHEN COALESCE(username, '') ILIKE $1 THEN 3
					WHEN COALESCE(username, '') ILIKE $1 || '%' THEN 2
					WHEN COALESCE(username, '') ILIKE '%' || $1 || '%' THEN 1
					ELSE 0
				END * 2
				+
				CASE WHEN COALESCE(first_name, '') ILIKE $1 THEN 3
					WHEN COALESCE(first_name, '') ILIKE $1 || '%' THEN 2
					WHEN COALESCE(first_name, '') ILIKE '%' || $1 || '%' THEN 1
					ELSE 0
				END * 1
				+
				CASE WHEN COALESCE(last_name, '') ILIKE $1 THEN 3
					WHEN COALESCE(last_name, '') ILIKE $1 || '%' THEN 2
					WHEN COALESCE(last_name, '') ILIKE '%' || $1 || '%' THEN 1
					ELSE 0
				END * 1
			) AS relevance
		FROM
			user_view
		WHERE
			COALESCE(email, '') ILIKE '%' || $1 || '%'
		OR	COALESCE(username, '') ILIKE '%' || $1 || '%'
		OR	COALESCE(first_name, '') ILIKE '%' || $1 || '%'
		OR	COALESCE(last_name, '') ILIKE '%' || $1 || '%'
		ORDER BY
			relevance DESC
		LIMIT 7
	`, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.UserSearchResult
	for rows.Next() {
		var ur model.UserSearchResult
		if err := rows.Scan(&ur.ID, &ur.Email, &ur.Username, &ur.FirstName, &ur.LastName, &ur.Relevance); err != nil {
			return nil, err
		}

		results = append(results, ur)
	}

	return results, nil
}

func (r *SearchRepository) SearchStockItems(
	ctx context.Context,
	exec db.PGExecutor,
	searchTerm string,
) ([]model.StockItemSearchResult, error) {

	rows, err := exec.Query(ctx, `
		SELECT
			stock_item_id,
			stock_code,
			description,
			(
				CASE WHEN stock_code ILIKE $1 THEN 3
					WHEN stock_code ILIKE $1 || '%' THEN 2
					WHEN stock_code ILIKE '%' || $1 || '%' THEN 1
					ELSE 0
				END * 3
				+
				CASE WHEN description ILIKE $1 THEN 3
					WHEN description ILIKE $1 || '%' THEN 2
					WHEN description ILIKE '%' || $1 || '%' THEN 1
					ELSE 0
				END * 2
			) AS relevance
		FROM
			stock_item
		WHERE
			stock_code ILIKE '%' || $1 || '%'
		OR	description ILIKE '%' || $1 || '%'
		ORDER BY
			relevance DESC
		LIMIT 7
	`, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.StockItemSearchResult
	for rows.Next() {
		var ur model.StockItemSearchResult
		if err := rows.Scan(&ur.StockItemID, &ur.StockCode, &ur.Description, &ur.Relevance); err != nil {
			return nil, err
		}

		results = append(results, ur)
	}

	return results, nil
}

func (r *SearchRepository) SearchResources(
	ctx context.Context,
	exec db.PGExecutor,
	searchTerm string,
) ([]model.ResourceSearchResult, error) {

	rows, err := exec.Query(ctx, `
		SELECT
			resource_id,
			type,
			reference,
			COALESCE(service_ownership_team_name, '') AS service_ownership_team_name,
			is_archived,
			(
				CASE WHEN reference ILIKE $1 THEN 3
					WHEN reference ILIKE $1 || '%' THEN 2
					WHEN reference ILIKE '%' || $1 || '%' THEN 1
					ELSE 0
				END * 3
				+
				CASE WHEN type ILIKE $1 THEN 3
					WHEN type ILIKE $1 || '%' THEN 2
					WHEN type ILIKE '%' || $1 || '%' THEN 1
					ELSE 0
				END * 2
				+
				CASE WHEN COALESCE(service_ownership_team_name, '') ILIKE $1 THEN 3
					WHEN COALESCE(service_ownership_team_name, '') ILIKE $1 || '%' THEN 2
					WHEN COALESCE(service_ownership_team_name, '') ILIKE '%' || $1 || '%' THEN 1
					ELSE 0
				END * 1
			) AS relevance
		FROM
			resource_view
		WHERE
			reference ILIKE '%' || $1 || '%'
		OR	type ILIKE '%' || $1 || '%'
		OR	COALESCE(service_ownership_team_name, '') ILIKE '%' || $1 || '%'
		ORDER BY
			relevance DESC
		LIMIT 7
	`, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.ResourceSearchResult
	for rows.Next() {
		var rr model.ResourceSearchResult
		if err := rows.Scan(
			&rr.ResourceID,
			&rr.Type,
			&rr.Reference,
			&rr.ServiceOwnershipTeamName,
			&rr.IsArchived,
			&rr.Relevance,
		); err != nil {
			return nil, err
		}

		results = append(results, rr)
	}

	return results, nil
}

func (r *SearchRepository) SearchServices(
	ctx context.Context,
	exec db.PGExecutor,
	searchTerm string,
) ([]model.ServiceSearchResult, error) {

	rows, err := exec.Query(ctx, `
		SELECT
			resource_service_id,
			resource_id,
			type,
			reference,
			status,
			started_at,
			COALESCE(started_by_username, '') AS started_by_username,
			(
				CASE WHEN reference ILIKE $1 THEN 3
					WHEN reference ILIKE $1 || '%' THEN 2
					WHEN reference ILIKE '%' || $1 || '%' THEN 1
					ELSE 0
				END * 3
				+
				CASE WHEN type ILIKE $1 THEN 3
					WHEN type ILIKE $1 || '%' THEN 2
					WHEN type ILIKE '%' || $1 || '%' THEN 1
					ELSE 0
				END * 2
				+
				CASE WHEN COALESCE(started_by_username, '') ILIKE $1 THEN 3
					WHEN COALESCE(started_by_username, '') ILIKE $1 || '%' THEN 2
					WHEN COALESCE(started_by_username, '') ILIKE '%' || $1 || '%' THEN 1
					ELSE 0
				END * 1
				+
				CASE WHEN COALESCE(notes, '') ILIKE $1 THEN 3
					WHEN COALESCE(notes, '') ILIKE $1 || '%' THEN 2
					WHEN COALESCE(notes, '') ILIKE '%' || $1 || '%' THEN 1
					ELSE 0
				END * 1
				+
				CASE WHEN status ILIKE $1 THEN 3
					WHEN status ILIKE $1 || '%' THEN 2
					WHEN status ILIKE '%' || $1 || '%' THEN 1
					ELSE 0
				END * 1
			) AS relevance
		FROM
			resource_service_view
		WHERE
			reference ILIKE '%' || $1 || '%'
		OR	type ILIKE '%' || $1 || '%'
		OR	COALESCE(started_by_username, '') ILIKE '%' || $1 || '%'
		OR	COALESCE(notes, '') ILIKE '%' || $1 || '%'
		OR	status ILIKE '%' || $1 || '%'
		ORDER BY
			relevance DESC
		LIMIT 7
	`, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.ServiceSearchResult
	for rows.Next() {
		var sr model.ServiceSearchResult
		if err := rows.Scan(
			&sr.ResourceServiceID,
			&sr.ResourceID,
			&sr.ResourceType,
			&sr.ResourceReference,
			&sr.Status,
			&sr.StartedAt,
			&sr.StartedByUsername,
			&sr.Relevance,
		); err != nil {
			return nil, err
		}

		results = append(results, sr)
	}

	return results, nil
}
