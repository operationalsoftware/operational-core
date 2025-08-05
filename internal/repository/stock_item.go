package repository

import (
	"app/internal/model"
	"app/pkg/db"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StockItemRepository struct{}

func NewStockItemRepository() *StockItemRepository {
	return &StockItemRepository{}
}

func (r *StockItemRepository) CreateStockItem(
	ctx context.Context,
	exec db.PGExecutor,
	stockItem *model.PostStockItem,
) error {

	insertUserStmt := `
INSERT INTO stock_item (
	stock_code,
	description
)
VALUES ($1, $2)
	`
	_, err := exec.Exec(
		ctx,

		insertUserStmt,
		stockItem.StockCode,
		stockItem.Description,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *StockItemRepository) UpdateStockItem(
	ctx context.Context,
	exec db.PGExecutor,
	stockCode string,
	input *model.PostStockItem,
) error {

	// get the user to check if it exists
	stockItem, err := r.GetStockItem(ctx, exec, stockCode)
	if err != nil {
		return err
	}
	if stockItem == nil {
		return fmt.Errorf("stock item does not exist")
	}

	query := `
UPDATE
	stock_item

SET
	stock_code = $2,
	description = $3

WHERE
	stock_code = $1
	`

	_, err = exec.Exec(
		ctx,
		query,

		stockCode,
		input.StockCode,
		input.Description,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *StockItemRepository) GetStockItem(
	ctx context.Context,
	exec db.PGExecutor,
	stockCode string,
) (*model.StockItem, error) {
	query := `
SELECT
    stock_code,
    description,
    created_at
FROM
    stock_item
WHERE
    stock_code = $1
	`

	stockItem := model.StockItem{}
	err := exec.QueryRow(ctx, query, stockCode).Scan(
		&stockItem.StockCode,
		&stockItem.Description,
		&stockItem.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &stockItem, nil
}

func (r *StockItemRepository) GetStockItems(
	ctx context.Context,
	exec db.PGExecutor,
	q *model.GetStockItemsQuery,
) ([]model.StockItem, error) {

	offset := (q.Page - 1) * q.PageSize
	limit := q.PageSize
	orderByClause, _ := q.Sort.ToOrderByClause(map[string]string{})

	if orderByClause == "" {
		orderByClause = "ORDER BY created_at DESC"
	}

	query := fmt.Sprintf(`
SELECT
    stock_code,
    description,
    created_at
FROM
    stock_item

%s

LIMIT $1 OFFSET $2
	`,
		orderByClause,
	)

	rows, err := exec.Query(ctx, query, limit, offset)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	stockItems := []model.StockItem{}
	for rows.Next() {
		var stockItem model.StockItem
		err := rows.Scan(
			&stockItem.StockCode,
			&stockItem.Description,
			&stockItem.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		stockItems = append(stockItems, stockItem)
	}

	return stockItems, nil
}

func (r *StockItemRepository) GetStockItemsCount(
	ctx context.Context,
	exec db.PGExecutor,
) (int, error) {

	query := `
SELECT
	COUNT(*)
FROM
	app_user
	`

	var count int
	err := exec.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *StockItemRepository) GetStockItemChanges(
	ctx context.Context,
	exec db.PGExecutor,
	stockCode string,
) ([]model.StockItemChange, error) {

	query := `
WITH RECURSIVE history_chain AS (
    SELECT *
    FROM stock_item_change
    WHERE stock_code = $1

    UNION ALL

    SELECT sic.*
    FROM stock_item_change sic
    INNER JOIN history_chain hc ON hc.stock_code_history = sic.stock_code
),
OldestChange AS (
    SELECT stock_code AS OldestChangeID
    FROM history_chain
    WHERE stock_code_history IS NULL
    LIMIT 1
)
SELECT
    sic.stock_code,
    sic.stock_code_history,
    sic.description,
    u.username AS changed_by_username,
    sic.changed_at,
    CASE
        WHEN sic.stock_code = OC.OldestChangeID THEN true
        ELSE false
    END AS IsCreation
FROM
    stock_item_change sic
LEFT JOIN app_user u ON sic.change_by = u.user_id
LEFT JOIN OldestChange OC ON TRUE
WHERE
    sic.stock_code = $1
ORDER BY
    sic.changed_at DESC;

`

	rows, err := exec.Query(ctx, query, stockCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var changes []model.StockItemChange
	for rows.Next() {
		var c model.StockItemChange
		err := rows.Scan(
			&c.StockCode,
			&c.StockCodeHistory,
			&c.Description,
			&c.ChangeByUsername,
			&c.ChangedAt,
			&c.IsCreation,
		)
		if err != nil {
			return nil, err
		}

		changes = append(changes, c)
	}

	return changes, nil
}

func (r *StockItemRepository) AddStockItemChange(
	ctx context.Context,
	exec db.PGExecutor,
	stockItemChange model.PostStockItemChange,
) error {

	insertQuery := `
INSERT INTO stock_item_change (
	stock_code,
	stock_code_history,
	description,
	change_by
)
VALUES ($1, $2, $3, $4)
	`
	_, err := exec.Exec(
		ctx,

		insertQuery,
		stockItemChange.StockCode,
		stockItemChange.StockCodeHistory,
		stockItemChange.Description,
		stockItemChange.ChangeBy,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *StockItemRepository) SearchStockItems(
	ctx context.Context,
	db *pgxpool.Pool,
	searchTerm string,
) ([]model.StockItemSearchResult, error) {

	rows, err := db.Query(ctx, `
		SELECT
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
		if err := rows.Scan(&ur.StockCode, &ur.Description, &ur.Relevance); err != nil {
			return nil, err
		}

		results = append(results, ur)
	}

	return results, nil
}
