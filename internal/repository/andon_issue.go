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
	severity,
	created_by
)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
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

// Create Group
func (r *AndonIssueRepository) CreateGroup(
	ctx context.Context,
	exec db.PGExecutor,
	andonIssueGroup model.NewAndonIssueGroup,
	userID int,
) error {

	if andonIssueGroup.ParentID != nil {
		isTopLevel, err := r.isTopLevelGroup(ctx, exec, *andonIssueGroup.ParentID)
		if err != nil {
			return fmt.Errorf("failed to validate parent group: %w", err)
		}
		if !isTopLevel {
			return fmt.Errorf("cannot create group deeper than 2 levels")
		}
	}

	query := `
INSERT INTO andon_issue (
	issue_name,
	parent_id,
	is_group,
	created_by
)
VALUES (
	$1,
	$2,
	$3,
	$4
)
`
	_, err := exec.Exec(
		ctx, query,

		andonIssueGroup.IssueName,
		andonIssueGroup.ParentID,
		true,
		userID,
	)

	return err
}

func (r *AndonIssueRepository) isTopLevelGroup(
	ctx context.Context,
	exec db.PGExecutor,
	groupID int,
) (bool, error) {

	var parentID *int
	query := `
SELECT
	parent_id
FROM
	andon_issue
WHERE
	andon_issue_id = $1 AND is_group = true
`
	err := exec.QueryRow(ctx, query, groupID).Scan(&parentID)
	if err != nil {
		return false, err
	}
	return parentID == nil, nil

}

var andonIssueCTE = `
RECURSIVE andon_issue_tree AS (
    SELECT
        ai.andon_issue_id,
        ai.issue_name,
        ai.parent_id,
        ARRAY[ai.issue_name] AS name_path,
        1 AS depth,
        ai.is_group,
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
        child.is_group,
        child.is_archived,
        (
            SELECT COUNT(*)
            FROM andon_issue c
            WHERE c.parent_id = child.andon_issue_id
        ) AS children_count
    FROM andon_issue child
    JOIN andon_issue_tree parent ON child.parent_id = parent.andon_issue_id
	WHERE parent.is_group = TRUE
),

final AS (
	SELECT
        ait.andon_issue_id,
        ait.issue_name,
        ait.parent_id,
        ait.name_path,
        ait.depth,
        ait.is_group,
        ait.is_archived,
        ait.children_count,
		ai.severity,
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
		LEFT JOIN team t ON
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

var andonIssueGroupSelect = `
SELECT
	andon_issue_id,
	issue_name,
	name_path,
	depth,
	parent_id,
	is_group,
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
	orderByClause, _ := q.Sort.ToOrderByClause(model.AndonIssue{})

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

func (r *AndonIssueRepository) ListIssuesAndGroups(
	ctx context.Context,
	exec db.PGExecutor,
	q model.ListAndonIssuesQuery,
) ([]model.AndonIssueNode, error) {

	whereClause := r.generateWhereClause(q)
	limit := q.PageSize
	offset := (q.Page - 1) * q.PageSize
	orderByClause, _ := q.Sort.ToOrderByClause(model.AndonIssueNode{})

	if orderByClause == "" {
		orderByClause = "ORDER BY created_at DESC"
	}

	query := `
WITH ` + andonIssueCTE + `

` + andonIssueGroupSelect + `

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

	var andonIssues []model.AndonIssueNode
	for rows.Next() {
		var andonIssue model.AndonIssueNode
		if err := rows.Scan(
			&andonIssue.AndonIssueID,
			&andonIssue.IssueName,
			&andonIssue.NamePath,
			&andonIssue.Depth,
			&andonIssue.ParentID,
			&andonIssue.IsGroup,
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

func (r *AndonIssueRepository) ListGroups(
	ctx context.Context,
	exec db.PGExecutor,
) ([]model.AndonIssueGroup, error) {

	query := `
WITH ` + andonIssueCTE + `

SELECT
	andon_issue_id,
	issue_name,
	name_path,
	parent_id
FROM
	final
WHERE
	is_group = true
`
	rows, err := exec.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var andonIssues []model.AndonIssueGroup
	for rows.Next() {
		var andonIssue model.AndonIssueGroup
		if err := rows.Scan(
			&andonIssue.AndonIssueID,
			&andonIssue.IssueName,
			&andonIssue.NamePath,
			&andonIssue.ParentID,
		); err != nil {
			return nil, err
		}

		andonIssues = append(andonIssues, andonIssue)
	}

	return andonIssues, nil
}

func (r *AndonIssueRepository) ListTopLevelGroups(
	ctx context.Context,
	exec db.PGExecutor,
) ([]model.AndonIssueGroup, error) {

	query := `
WITH ` + andonIssueCTE + `

SELECT
	andon_issue_id,
	issue_name,
	name_path,
	parent_id
FROM
	final
WHERE
	is_group = true AND parent_id IS NULL
`
	rows, err := exec.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var andonIssues []model.AndonIssueGroup
	for rows.Next() {
		var andonIssue model.AndonIssueGroup
		if err := rows.Scan(
			&andonIssue.AndonIssueID,
			&andonIssue.IssueName,
			&andonIssue.NamePath,
			&andonIssue.ParentID,
		); err != nil {
			return nil, err
		}

		andonIssues = append(andonIssues, andonIssue)
	}

	return andonIssues, nil
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

	var existingAssignedToTeam *int
	// Ensure exists
	existing, err := r.GetByID(ctx, exec, andonIssueID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("andon issue with ID %d not found", andonIssueID)
	}
	existingAssignedToTeam = existing.AssignedToTeam

	severityChanged := false
	if existing.Severity == nil || update.Severity != *existing.Severity {
		severityChanged = true
	}

	// Check if anything changed
	hasChange := update.IssueName != existing.IssueName ||
		update.ParentID != existing.ParentID ||
		update.IsArchived != existing.IsArchived ||
		update.AssignedToTeam != existingAssignedToTeam ||
		severityChanged

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
		"issue_name":       update.IssueName,
		"parent_id":        update.ParentID,
		"is_archived":      update.IsArchived,
		"assigned_to_team": update.AssignedToTeam,
		"severity":         update.Severity,
		"updated_by":       userID,
		"andon_issue_id":   andonIssueID,
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

func (r *AndonIssueRepository) GetIssueHierarchy(
	ctx context.Context,
	exec db.PGExecutor,
	issueID int,
) ([]int, error) {

	query := `
WITH RECURSIVE issue_hierarchy AS (
    SELECT
        andon_issue_id,
        issue_name,
        parent_id
    FROM andon_issue
    WHERE andon_issue_id = $1  -- pass the issue ID here

    UNION ALL

    SELECT
        parent.andon_issue_id,
        parent.issue_name,
        parent.parent_id
    FROM andon_issue parent
    JOIN issue_hierarchy child ON child.parent_id = parent.andon_issue_id
)
SELECT andon_issue_id
FROM issue_hierarchy
ORDER BY parent_id NULLS FIRST;
`

	rows, err := exec.Query(ctx, query, issueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []int
	for rows.Next() {
		var name int
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return names, nil
}
