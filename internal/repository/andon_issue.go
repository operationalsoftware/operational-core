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
	userID int,
) error {

	query := `
INSERT INTO andon_issue (
	issue_name,
	parent_id,
	assigned_to_team,
	severity
	created_by
)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6
)
`
	_, err := exec.Exec(
		ctx, query,

		andonIssue.IssueName,
		andonIssue.ParentID,
		andonIssue.AssignedToTeam,
		andonIssue.Severity,
		userID,
	)

	return err
}

var andonIssueCTE = `
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
        ) AS children_count
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
        ) AS children_count
    FROM andon_issue child
    JOIN andon_issue_tree parent ON child.parent_id = parent.andon_issue_id
),

final AS (
	SELECT
        ait.andon_issue_id,
        ait.issue_name,
        ait.parent_id,
        ait.name_path,
        ait.depth,
        ait.is_archived,
        ait.children_count,
		ai.resolvable_by_raiser,
		ai.will_stop_process,
		ai.assigned_to_team,
		t.team_name AS assigned_to_team_name,
		ai.created_at,
		ai.created_by,
		cu.username AS created_by_username,
		ai.updated_at,
		ai.updated_by,
		uu.username AS updated_by_username
	
	FROM
		andon_issue_tree ait
		INNER JOIN andon_issue ai USING(andon_issue_id)
		INNER JOIN team t ON
			t.team_id = ai.assigned_to_team
		INNER JOIN app_user cu ON
			cu.user_id = ai.created_by
		LEFT JOIN app_user uu ON
			uu.user_id = ai.updated_by
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
	severity,
	assigned_to_team,
	assigned_to_team_name,

	created_at,
	created_by,
	created_by_username,
	updated_at,
	updated_by,
	updated_by_username
`

// Read
func (r *AndonIssueRepository) GetByID(
	ctx context.Context,
	exec db.PGExecutor,
	andonIssueID int,
) (*model.AndonIssue, error) {
	query := `
WITH ` + andonIssueCTE + `

` + andonIssueSelect + `

FROM
	final

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
		&andonIssue.Severity,
		&andonIssue.AssignedToTeam,
		&andonIssue.AssignedToTeamName,
		&andonIssue.CreatedAt,
		&andonIssue.CreatedBy,
		&andonIssue.CreatedByUsername,
		&andonIssue.UpdatedAt,
		&andonIssue.UpdatedBy,
		&andonIssue.UpdatedByUsername,
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
WITH ` + andonIssueCTE + `

` + andonIssueSelect + `

FROM
	final

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
			&andonIssue.Severity,
			&andonIssue.AssignedToTeam,
			&andonIssue.AssignedToTeamName,
			&andonIssue.CreatedAt,
			&andonIssue.CreatedBy,
			&andonIssue.CreatedByUsername,
			&andonIssue.UpdatedAt,
			&andonIssue.UpdatedBy,
			&andonIssue.UpdatedByUsername,
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
WITH ` + andonIssueCTE + `

SELECT
	COUNT(*)

FROM
	final

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
	userID int,
) error {

	// Ensure exists
	existing, err := r.GetByID(ctx, exec, andonIssueID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("andon issue with ID %d not found", andonIssueID)
	}

	// Check if anything changed
	hasChange := update.IssueName != existing.IssueName ||
		update.ParentID != existing.ParentID ||
		update.IsArchived != existing.IsArchived ||
		update.AssignedToTeam != existing.AssignedToTeam ||
		update.Severity != existing.Severity

	if !hasChange {
		return nil
	}

	namedQuery := `
UPDATE andon_issue
SET
	issue_name = :issue_name,
	parent_id = :parent_id,
	is_archived = :is_archived,
	assigned_to_team = :assigned_to_team,
	severity = :severity,
	updated_by = :updated_by,
	updated_at = NOW()
WHERE
	andon_issue_id = :andon_issue_id
`

	query, params, err := db.BindNamed(namedQuery, map[string]any{
		"issue_name":      update.IssueName,
		"parent_id":       update.ParentID,
		"is_archived":     update.IsArchived,
		"assigned_to_tem": update.AssignedToTeam,
		"severity":        update.Severity,
		"updated_by":      userID,
		"andon_issue_id":  andonIssueID,
	})
	if err != nil {
		return err
	}

	_, err = exec.Exec(
		ctx,
		query,
		params...,
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
            FROM andon_issue ai
            WHERE ai.parent_id = $1

            UNION ALL

            SELECT
                child.andon_issue_id,
                child.parent_id,
                child.is_archived
            FROM andon_issue child
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
