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
	assigned_team,
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
		andonIssue.AssignedTeam,
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
	assigned_team,
	assigned_team_name,

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
) (*model.AndonIssueNode, error) {
	query := `
SELECT
	andon_issue_id,
	issue_name,
	name_path,
	depth,
	down_depth,
	parent_id,
	is_archived,
	is_group,
	children_count,
	severity,
	assigned_team,
	assigned_team_name,

	created_at,
	created_by,
	created_by_username,
	updated_at,
	updated_by,
	updated_by_username
FROM
	andon_issue_tree_view

WHERE
	andon_issue_id = $1
`

	var andonIssue model.AndonIssueNode
	err := exec.QueryRow(ctx, query, andonIssueID).Scan(
		&andonIssue.AndonIssueID,
		&andonIssue.IssueName,
		&andonIssue.NamePath,
		&andonIssue.Depth,
		&andonIssue.DownDepth,
		&andonIssue.ParentID,
		&andonIssue.IsArchived,
		&andonIssue.IsGroup,
		&andonIssue.ChildrenCount,
		&andonIssue.Severity,
		&andonIssue.AssignedTeam,
		&andonIssue.AssignedTeamName,
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

func (r *AndonIssueRepository) GetIssueByID(
	ctx context.Context,
	exec db.PGExecutor,
	andonIssueID int,
) (*model.AndonIssue, error) {
	query := `
SELECT
	andon_issue_id,
	issue_name,
	name_path,
	depth,
	parent_id,
	is_archived,
	children_count,
	severity,
	assigned_team,
	assigned_team_name,

	created_at,
	created_by,
	created_by_username,
	updated_at,
	updated_by,
	updated_by_username
FROM
	andon_issue_tree_view

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
		&andonIssue.AssignedTeam,
		&andonIssue.AssignedTeamName,
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

func (r *AndonIssueRepository) GetGroupByID(
	ctx context.Context,
	exec db.PGExecutor,
	andonIssueGroupID int,
) (*model.AndonIssueGroup, error) {
	query := `
SELECT
	andon_issue_id,
	issue_name,
	name_path,
	depth,
	down_depth,
	parent_id
FROM
	andon_issue_group_view

WHERE
	andon_issue_id = $1
`

	var andonIssue model.AndonIssueGroup
	err := exec.QueryRow(ctx, query, andonIssueGroupID).Scan(
		&andonIssue.AndonIssueID,
		&andonIssue.IssueName,
		&andonIssue.NamePath,
		&andonIssue.Depth,
		&andonIssue.DownDepth,
		&andonIssue.ParentID,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &andonIssue, nil
}

func (r *AndonIssueRepository) ListIssues(
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
SELECT
	andon_issue_id,
	issue_name,
	name_path,
	depth,
	parent_id,
	is_archived,
	children_count,
	severity,
	assigned_team,
	assigned_team_name,

	created_at,
	created_by,
	created_by_username,
	updated_at,
	updated_by,
	updated_by_username

FROM
	andon_issue_view

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
			&andonIssue.AssignedTeam,
			&andonIssue.AssignedTeamName,
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

	query := andonIssueGroupSelect + `

FROM
	andon_issue_tree_view

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
			&andonIssue.AssignedTeam,
			&andonIssue.AssignedTeamName,
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
SELECT
	COUNT(*)

FROM
	andon_issue_tree_view

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
SELECT
	andon_issue_id,
	issue_name,
	name_path,
	parent_id,
	depth,
	down_depth
FROM
	andon_issue_group_view
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
			&andonIssue.Depth,
			&andonIssue.DownDepth,
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
SELECT
	andon_issue_id,
	issue_name,
	name_path,
	parent_id,
	down_depth
FROM
	andon_issue_group_view
WHERE
	parent_id IS NULL
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
			&andonIssue.DownDepth,
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

	// Ensure exists
	existing, err := r.GetIssueByID(ctx, exec, andonIssueID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("andon issue with ID %d not found", andonIssueID)
	}

	severityChanged := false
	if existing.Severity == "" || update.Severity != existing.Severity {
		severityChanged = true
	}

	// Check if anything changed
	hasChange := update.IssueName != existing.IssueName ||
		update.ParentID != existing.ParentID ||
		update.IsArchived != existing.IsArchived ||
		update.AssignedTeam != existing.AssignedTeam ||
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
	assigned_team = :assigned_team,
	severity = :severity,
	updated_by = :updated_by,
	updated_at = NOW()
WHERE
	andon_issue_id = :andon_issue_id
`

	query, params, err := db.BindNamed(namedQuery, map[string]any{
		"issue_name":     update.IssueName,
		"parent_id":      update.ParentID,
		"is_archived":    update.IsArchived,
		"assigned_team":  update.AssignedTeam,
		"severity":       update.Severity,
		"updated_by":     userID,
		"andon_issue_id": andonIssueID,
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

// Update issue group
func (r *AndonIssueRepository) UpdateGroup(
	ctx context.Context,
	exec db.PGExecutor,
	andonIssueID int,
	update model.AndonIssueGroupUpdate,
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
		update.ParentID != existing.ParentID

	if !hasChange {
		return nil
	}

	namedQuery := `
UPDATE andon_issue
SET
	issue_name = :issue_name,
	parent_id = :parent_id,
	updated_by = :updated_by,
	updated_at = NOW()
WHERE
	andon_issue_id = :andon_issue_id
`

	query, params, err := db.BindNamed(namedQuery, map[string]any{
		"issue_name":     update.IssueName,
		"parent_id":      update.ParentID,
		"updated_by":     userID,
		"andon_issue_id": andonIssueID,
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

// Get list of parent ids of an andon issue in sorted order
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
