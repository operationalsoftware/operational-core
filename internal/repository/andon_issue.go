package repository

import (
	"app/internal/model"
	"app/pkg/db"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type AndonIssueRepository struct{}

func NewAndonIssueRepository() *AndonIssueRepository {
	return &AndonIssueRepository{}
}

// Create
func (r *AndonIssueRepository) Create(
	ctx context.Context,
	exec db.PGExecutor,
	andonIssue model.NewAndonIssue,
) error {

	query := `
INSERT INTO andon_issue (
	issue_name,
	parent_id
)
VALUES (
	$1,
	$2
)
`
	_, err := exec.Exec(
		ctx, query,

		andonIssue.IssueName,
		andonIssue.ParentID,
	)

	return err
}

var andonIssueTreeCTE = `
RECURSIVE andon_issue_tree AS (
    SELECT
        ai.andon_issue_id,
        ai.issue_name,
        ai.parent_id,
        ARRAY[ai.issue_name] AS name_path,
        1 AS depth,
        ai.is_archived,
        (
            SELECT COUNT(*)
            FROM andon_issue c
            WHERE c.parent_id = ai.andon_issue_id
        ) AS children_count, 
        ai.created_at,
        ai.created_by,
        ai.updated_at,
        ai.updated_by
    FROM andon_issue ai
    WHERE ai.parent_id IS NULL

    UNION ALL

    SELECT
        child.andon_issue_id,
        child.issue_name,
        child.parent_id,
        parent.name_path || child.issue_name,
        parent.depth + 1,
        child.is_archived,
        (
            SELECT COUNT(*)
            FROM andon_issue c
            WHERE c.parent_id = child.andon_issue_id
        ) AS children_count,
        child.created_at,
        child.created_by,
        child.updated_at,
        child.updated_by
    FROM andon_issue child
    JOIN andon_issue_tree parent ON child.parent_id = parent.andon_issue_id
)
`

var andonIssueSelect = `
SELECT 
	andon_issue_id, 
	issue_name, 
	name_path, 
	depth, 
	parent_id, 
	is_archived, 
	children_count,
	created_at, 
	created_by, 
	updated_at, 
	updated_by
`

// Read
func (r *AndonIssueRepository) GetByID(
	ctx context.Context,
	exec db.PGExecutor,
	andonIssueID int,
) (*model.AndonIssue, error) {
	query := `
WITH ` + andonIssueTreeCTE + `

` + andonIssueSelect + `

FROM
	andon_issue_tree

WHERE
	andon_issue_id = $1
`

	var andonIssue model.AndonIssue
	err := exec.QueryRow(ctx, query, andonIssueID).Scan(
		&andonIssue.AndonIssueID,
		&andonIssue.IssueName,
		&andonIssue.NamePath,
		&andonIssue.Depth,
		&andonIssue.ParentID,
		&andonIssue.IsArchived,
		&andonIssue.ChildrenCount,
		&andonIssue.CreatedAt,
		&andonIssue.CreatedBy,
		&andonIssue.UpdatedAt,
		&andonIssue.UpdatedBy,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &andonIssue, nil
}

func (r *AndonIssueRepository) List(
	ctx context.Context,
	exec db.PGExecutor,
	q model.ListAndonIssuesQuery,
) ([]model.AndonIssue, error) {

	whereClause := r.generateWhereClause(q)
	limit := q.PageSize
	offset := (q.Page - 1) * q.PageSize
	orderByClause, err := q.Sort.ToOrderByClause(model.AndonIssue{})

	if orderByClause == "" {
		orderByClause = "ORDER BY name_path ASC"
	}

	query := `
WITH ` + andonIssueTreeCTE + `

` + andonIssueSelect + `

FROM
	andon_issue_tree

` + whereClause + `

` + orderByClause + `

LIMIT $1 OFFSET $2
`
	rows, err := exec.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var andonIssues []model.AndonIssue
	for rows.Next() {
		var andonIssue model.AndonIssue
		if err := rows.Scan(
			&andonIssue.AndonIssueID,
			&andonIssue.IssueName,
			&andonIssue.NamePath,
			&andonIssue.Depth,
			&andonIssue.ParentID,
			&andonIssue.IsArchived,
			&andonIssue.ChildrenCount,
			&andonIssue.CreatedAt,
			&andonIssue.CreatedBy,
			&andonIssue.UpdatedAt,
			&andonIssue.UpdatedBy,
		); err != nil {
			return nil, err
		}

		andonIssues = append(andonIssues, andonIssue)
	}

	return andonIssues, nil
}

func (r *AndonIssueRepository) Count(
	ctx context.Context,
	exec db.PGExecutor,
	q model.ListAndonIssuesQuery,
) (int, error) {

	whereClause := r.generateWhereClause(q)

	query := `
WITH ` + andonIssueTreeCTE + `

SELECT
	COUNT(*)

FROM
	andon_issue_tree

` + whereClause

	var count int
	err := exec.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *AndonIssueRepository) generateWhereClause(q model.ListAndonIssuesQuery) string {
	if !q.ShowArchived {
		return `
WHERE
	is_archived = false
`
	} else {
		return ""
	}
}

// Update
func (r *AndonIssueRepository) Update(
	ctx context.Context,
	exec db.PGExecutor,
	andonIssueID int,
	update model.AndonIssueUpdate,
) error {

	// Ensure exists
	existing, err := r.GetByID(ctx, exec, andonIssueID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("andon issue with ID %d not found", andonIssueID)
	}

	query := `
UPDATE team
SET
	issue_name = $1,
	parent_id = $2,
    is_archived = $3
WHERE
	andon_issue_id = $4
`
	_, err = exec.Exec(
		ctx,
		query,

		update.IssueName,
		update.ParentID,
		update.IsArchived,
		andonIssueID,
	)

	return err
}

func (r *AndonIssueRepository) HasActiveChildIssues(
	ctx context.Context,
	exec db.PGExecutor,
	andonIssueID int,
) (bool, error) {
	var hasActive bool

	err := exec.QueryRow(ctx, `
        WITH RECURSIVE child_issues AS (
            SELECT
                ai.andon_issue_id,
                ai.parent_id,
                ai.is_archived
            FROM andon_issues ai
            WHERE ai.parent_id = $1

            UNION ALL

            SELECT
                child.andon_issue_id,
                child.parent_id,
                child.is_archived
            FROM andon_issues child
            JOIN child_issues parent ON child.parent_id = parent.andon_issue_id
        )
        SELECT EXISTS (
            SELECT 1
            FROM child_issues
            WHERE is_archived = false
        );
    `, andonIssueID).Scan(&hasActive)

	if err != nil {
		return false, err
	}

	return hasActive, nil
}
