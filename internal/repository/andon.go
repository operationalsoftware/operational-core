package repository

import (
	"app/internal/components"
	"app/internal/model"
	"app/pkg/db"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

type AndonRepository struct{}

func NewAndonRepository() *AndonRepository {
	return &AndonRepository{}
}

func (r *AndonRepository) CreateAndonEvent(
	ctx context.Context,
	exec db.PGExecutor,
	andonEvent model.NewAndonEvent,
	userID int,
) error {

	insertQuery := `
INSERT INTO
  andon_change (
    change_by,
    andon_event_id,
    raised_by,
    status
)
VALUES (
    $1, $2, $1, 'Outstanding'
)
`

	query := `
INSERT INTO andon_event (
	issue_description,
	issue_id,
	source,
	location,
	status,
    linked_entity_id,
    linked_entity_type,
	raised_by
)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8
)
RETURNING andon_event_id
`

	var newAndonID int
	err := exec.QueryRow(
		ctx, query,

		andonEvent.IssueDescription,
		andonEvent.IssueID,
		andonEvent.Source,
		andonEvent.Location,
		"Outstanding",
		andonEvent.LinkedEntityID,
		andonEvent.LinkedEntityType,
		userID,
	).Scan(&newAndonID)
	if err != nil {
		return err
	}

	_, err = exec.Exec(ctx, insertQuery, userID, newAndonID)
	if err != nil {
		return err
	}

	return nil
}

func (r *AndonRepository) GetAndonEventByID(
	ctx context.Context,
	exec db.PGExecutor,
	andonEventID int,
	userID int,
) (*model.AndonEvent, error) {

	var userTeamIDs []int
	err := exec.QueryRow(ctx, `
SELECT
	array_agg(team_id)
FROM
	user_team
WHERE
	user_id = $1
`, userID).Scan(&userTeamIDs)
	if err != nil {
		return nil, err
	}

	query := `
SELECT
	issue_description,
	issue_id,
	source,
	location,
	status,
	assigned_team,
	assigned_team_name,
	raised_by_username,
	raised_at,
	acknowledged_by_username,
	acknowledged_at,
	resolved_by_username,
	resolved_at,
	cancelled_by_username,
	cancelled_at,
	last_updated,
	name_path,
	severity,
	(
		severity = 'Info'
		OR
		(
			severity = 'Requires Intervention'
			AND 
			assigned_team IN (SELECT team_id FROM user_team WHERE user_id = $2)
		)
	) AS can_user_acknowledge,
	(
		(
			severity = 'Self-resolvable'
			OR
			severity = 'Requires Intervention'
		)
		AND 
		assigned_team IN (
			SELECT team_id FROM user_team WHERE user_id = $2
		)
	) AS can_user_resolve,
	(
		status <> 'Resolved'
		AND (
			(
				severity IN ('Info', 'Self-resolvable')
				AND (
					raised_by = $2
					OR assigned_team IN (SELECT team_id FROM user_team WHERE user_id = $2)
				)
			)
			OR (
				severity = 'Requires Intervention'
				AND assigned_team IN (SELECT team_id FROM user_team WHERE user_id = $2)
			)
		)
	) AS can_user_cancel
FROM
	andon_view
WHERE
	andon_event_id = $1
`

	var andonEvent model.AndonEvent
	err = exec.QueryRow(
		ctx, query, andonEventID, userID,
	).Scan(
		&andonEvent.IssueDescription,
		&andonEvent.IssueID,
		&andonEvent.Source,
		&andonEvent.Location,
		&andonEvent.Status,
		&andonEvent.AssignedTeam,
		&andonEvent.AssignedTeamName,
		&andonEvent.RaisedByUsername,
		&andonEvent.RaisedAt,
		&andonEvent.AcknowledgedByUsername,
		&andonEvent.AcknowledgedAt,
		&andonEvent.ResolvedByUsername,
		&andonEvent.ResolvedAt,
		&andonEvent.CancelledByUsername,
		&andonEvent.CancelledAt,
		&andonEvent.LastUpdated,
		&andonEvent.NamePath,
		&andonEvent.Severity,
		&andonEvent.CanUserAcknowledge,
		&andonEvent.CanUserResolve,
		&andonEvent.CanUserCancel,
	)
	if err != nil {
		return nil, err
	}

	return &andonEvent, err
}

func (r *AndonRepository) ListAndons(
	ctx context.Context,
	exec db.PGExecutor,
	q model.ListAndonQuery,
	currentUserID int,
) ([]model.AndonEvent, error) {

	var userTeamIDs []int
	err := exec.QueryRow(ctx, `
SELECT
	array_agg(team_id)
FROM
	user_team
WHERE
	user_id = $1
`, currentUserID).Scan(&userTeamIDs)
	if err != nil {
		return nil, err
	}

	whereClause, args := generateWhereClause(q)

	currentUserIDPlaceholder := fmt.Sprintf("$%d", len(args)+1)
	limitPlaceholder := fmt.Sprintf("$%d", len(args)+2)
	offsetPlaceholder := fmt.Sprintf("$%d", len(args)+3)

	query := `
SELECT
	andon_event_id,
	issue_description,
	issue_id,
	source,
	location,
	status,
	raised_by_username,
	raised_at,
	acknowledged_by_username,
	acknowledged_at,
	resolved_by_username,
	resolved_at,
	last_updated,
	assigned_team,
	assigned_team_name,
	issue_name,
	name_path,
	severity,
	(
	severity = 'Info'
	OR (
		severity = 'Requires Intervention'
		AND assigned_team IN (
			SELECT team_id FROM user_team WHERE user_id = ` + currentUserIDPlaceholder + `
		)
	)
	) AS can_user_acknowledge,
	(
		(severity = 'Self-resolvable' OR severity = 'Requires Intervention')
		AND assigned_team IN (
			SELECT team_id FROM user_team WHERE user_id = ` + currentUserIDPlaceholder + `
		)
	) AS can_user_resolve,
	(
		status <> 'Resolved'
		AND (
			(
				severity IN ('Info', 'Self-resolvable')
				AND (
					raised_by = ` + currentUserIDPlaceholder + `
					OR assigned_team IN (SELECT team_id FROM user_team WHERE user_id = ` + currentUserIDPlaceholder + `)
				)
			)
			OR (
				severity = 'Requires Intervention'
				AND assigned_team IN (SELECT team_id FROM user_team WHERE user_id = ` + currentUserIDPlaceholder + `)
			)
		)
	) AS can_user_cancel
FROM andon_view
`

	limit := q.PageSize
	offset := (q.Page - 1) * q.PageSize
	orderByClause, _ := q.Sort.ToOrderByClause(model.AndonEvent{})
	defaultOrderBy := "resolved_at"
	defaultOrderByDirection := "asc"

	if q.OrderBy != "" {
		defaultOrderBy = q.OrderBy
	}
	if q.OrderByDirection != "" {
		defaultOrderByDirection = q.OrderByDirection
	}

	if orderByClause == "" {
		orderByClause = fmt.Sprintf("ORDER BY %s %s", defaultOrderBy, defaultOrderByDirection)
	}

	finalQuery := query + "\n" + whereClause + "\n" + orderByClause + "\n" +
		fmt.Sprintf("LIMIT %s OFFSET %s", limitPlaceholder, offsetPlaceholder)

	rows, err := exec.Query(ctx, finalQuery, append(args, currentUserID, limit, offset)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.AndonEvent

	for rows.Next() {
		var event model.AndonEvent

		if err := rows.Scan(
			&event.AndonEventID,
			&event.IssueDescription,
			&event.IssueID,
			&event.Source,
			&event.Location,
			&event.Status,
			&event.RaisedByUsername,
			&event.RaisedAt,
			&event.AcknowledgedByUsername,
			&event.AcknowledgedAt,
			&event.ResolvedByUsername,
			&event.ResolvedAt,
			&event.LastUpdated,
			&event.AssignedTeam,
			&event.AssignedTeamName,
			&event.IssueName,
			&event.NamePath,
			&event.Severity,
			&event.CanUserAcknowledge,
			&event.CanUserResolve,
			&event.CanUserCancel,
		); err != nil {
			return nil, err
		}

		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *AndonRepository) Count(
	ctx context.Context,
	exec db.PGExecutor,
	q model.ListAndonQuery,
) (int, error) {

	query := `
SELECT COUNT(*) FROM andon_view
`

	whereClause, args := generateWhereClause(q)
	finalQuery := query + "\n" + whereClause

	var count int
	err := exec.QueryRow(ctx, finalQuery, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *AndonRepository) GetAvailableFilters(
	ctx context.Context,
	exec db.PGExecutor,
	baseFilters model.AndonFilters,
) (model.AndonAvailableFilters, error) {

	mapping := map[string]string{
		"IssueIn":                  "issue_name",
		"SeverityIn":               "severity",
		"TeamIn":                   "assigned_team_name",
		"LocationIn":               "location",
		"StatusIn":                 "status",
		"RaisedByUsernameIn":       "raised_by_username",
		"AcknowledgedByUsernameIn": "acknowledged_by_username",
		"ResolvedByUsernameIn":     "resolved_by_username",
	}

	avail := model.AndonAvailableFilters{}

	// helper to select into a *[]string
	var collect = func(key string, dest *[]string) error {
		queryFilters := model.ListAndonQuery{
			StartDate:              baseFilters.StartDate,
			EndDate:                baseFilters.EndDate,
			Issues:                 baseFilters.Issues,
			Serverities:            baseFilters.Severities,
			Teams:                  baseFilters.Teams,
			Locations:              baseFilters.Locations,
			Statuses:               baseFilters.Statuses,
			RaisedByUsername:       baseFilters.RaisedByUsername,
			AcknowledgedByUsername: baseFilters.AcknowledgedByUsername,
			ResolvedByUsername:     baseFilters.ResolvedByUsername,
		}

		switch key {
		case "IssueIn":
			queryFilters.Issues = nil
		case "SeverityIn":
			queryFilters.Serverities = nil
		case "TeamIn":
			queryFilters.Teams = nil
		case "LocationIn":
			queryFilters.Locations = nil
		case "StatusIn":
			queryFilters.Statuses = nil
		case "RaisedByUsernameIn":
			queryFilters.RaisedByUsername = nil
		case "AcknowledgedByUsernameIn":
			queryFilters.AcknowledgedByUsername = nil
		case "ResolvedByUsernameIn":
			queryFilters.ResolvedByUsername = nil
		}

		where, args := generateWhereClause(queryFilters)
		col := mapping[key]

		if col != "" {
			if where == "" {
				where = "WHERE " + col + " IS NOT NULL"
			} else {
				where += "\nAND " + col + " IS NOT NULL"
			}
		}

		query := `
SELECT DISTINCT ` + col + ` AS val
FROM andon_view
` + where + `
ORDER BY val ASC
`

		rows, err := exec.Query(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var v pgtype.Text
			if err := rows.Scan(&v); err != nil {
				return err
			}
			*dest = append(*dest, v.String)
		}

		return rows.Err()
	}

	if err := collect("IssueIn", &avail.IssueIn); err != nil {
		return avail, err
	}
	if err := collect("SeverityIn", &avail.SeverityIn); err != nil {
		return avail, err
	}
	if err := collect("TeamIn", &avail.TeamIn); err != nil {
		return avail, err
	}
	if err := collect("LocationIn", &avail.LocationIn); err != nil {
		return avail, err
	}
	if err := collect("StatusIn", &avail.StatusIn); err != nil {
		return avail, err
	}
	if err := collect("RaisedByUsernameIn", &avail.RaisedByUsernameIn); err != nil {
		return avail, err
	}
	if err := collect("AcknowledgedByUsernameIn", &avail.AcknowledgedByUsernameIn); err != nil {
		return avail, err
	}
	if err := collect("ResolvedByUsernameIn", &avail.ResolvedByUsernameIn); err != nil {
		return avail, err
	}

	return avail, nil
}

func (r *AndonRepository) GetAndonChangeLog(
	ctx context.Context,
	exec db.PGExecutor,
	andonEventID int,
) ([]model.AndonChange, error) {

	query := `
SELECT
	ac.andon_event_id,
	ac.issue_description,
	ac.issue_id,
	ac.location,
	ac.status,
	raised_user.username AS raised_by_username,
	ac.raised_at,
	au.username AS acknowledged_by_username,
	ac.acknowledged_at,
	ru.username AS resolved_by_username,
	ac.resolved_at,
	cu.username AS cancelled_by_username,
	ac.cancelled_at,
	change_user.username AS change_by_username,
	ac.change_at,
    CASE
        WHEN ac.change_at = MIN(ac.change_at) OVER (PARTITION BY ac.andon_event_id)
        THEN true
        ELSE false
    END AS IsCreation
FROM
	andon_change AS ac
LEFT JOIN
	app_user AS raised_user ON ac.raised_by = raised_user.user_id
LEFT JOIN
	app_user AS au ON ac.acknowledged_by = au.user_id
LEFT JOIN
	app_user AS ru ON ac.resolved_by = ru.user_id
LEFT JOIN
	app_user AS cu ON ac.cancelled_by = cu.user_id
INNER JOIN
	app_user AS change_user ON ac.change_by = change_user.user_id
LEFT JOIN
	andon_issue ai ON ac.issue_id = ai.andon_issue_id

WHERE andon_event_id = $1
ORDER BY ac.change_at DESC
`

	rows, err := exec.Query(ctx, query, andonEventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.AndonChange
	for rows.Next() {
		var event model.AndonChange
		if err := rows.Scan(
			&event.AndonEventID,
			&event.IssueDescription,
			&event.IssueID,
			&event.Location,
			&event.Status,
			&event.RaisedByUsername,
			&event.RaisedAt,
			&event.AcknowledgedByUsername,
			&event.AcknowledgedAt,
			&event.ResolvedByUsername,
			&event.ResolvedAt,
			&event.CancelledByUsername,
			&event.CancelledAt,
			&event.ChangeByUsername,
			&event.ChangeAt,
			&event.IsCreation,
		); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *AndonRepository) GetAndonComments(
	ctx context.Context,
	exec db.PGExecutor,
	andonEventID int,
) ([]components.Comment, error) {

	query := `
SELECT
	c.comment_id,
	c.entity_id,
	c.comment,
	u.username,
	c.commented_at
FROM comment c
LEFT JOIN
	app_user AS u ON c.commented_by = u.user_id
WHERE c.entity = 'andon' AND c.entity_id = $1
ORDER BY c.commented_at ASC
`

	rows, err := exec.Query(ctx, query, andonEventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []components.Comment
	for rows.Next() {
		var comment components.Comment
		if err := rows.Scan(
			&comment.CommentID, // maps to comment_id
			&comment.EntityID,  // maps to entity_id
			&comment.Comment,
			&comment.CommentedByUsername,
			&comment.CommentedAt,
		); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *AndonRepository) AcknowledgeAndonEvent(
	ctx context.Context,
	exec db.PGExecutor,
	andonEventID int,
	userID int,
) error {

	insertQuery := `
INSERT INTO 
  andon_change (
    change_by, 
    andon_event_id, 
    acknowledged_by, 
    acknowledged_at, 
    status
) 
VALUES (
    $1, $2, $1, NOW(), 'Acknowledged'
)
`

	updateQuery := `
UPDATE 
  andon_event
SET 
  acknowledged_by = $1,
  acknowledged_at = NOW(),
  status = 'Acknowledged',
  last_updated = NOW()
WHERE 
  andon_event_id = $2
`

	_, err := exec.Exec(ctx, insertQuery, userID, andonEventID)
	if err != nil {
		return err
	}

	_, err = exec.Exec(ctx, updateQuery, userID, andonEventID)
	if err != nil {
		return err
	}

	return nil
}

func (r *AndonRepository) ResolveAndonEvent(
	ctx context.Context,
	exec db.PGExecutor,
	andonEventID int,
	userID int,
) error {

	insertQuery := `
INSERT INTO 
  andon_change (
    change_by, 
    andon_event_id, 
    resolved_by, 
    resolved_at, 
    status
) 
VALUES (
    $1, $2, $1, NOW(), 'Resolved'
)
`

	updateQuery := `
UPDATE 
  andon_event
SET 
  resolved_by = $1,
  resolved_at = NOW(),
  status = 'Resolved',
  last_updated = NOW()
WHERE 
  andon_event_id = $2
`

	_, err := exec.Exec(ctx, insertQuery, userID, andonEventID)
	if err != nil {
		return err
	}

	_, err = exec.Exec(ctx, updateQuery, userID, andonEventID)
	if err != nil {
		return err
	}

	return nil
}

func (r *AndonRepository) CancelAndonEvent(
	ctx context.Context,
	exec db.PGExecutor,
	andonEventID int,
	userID int,
) error {

	insertQuery := `
INSERT INTO 
  andon_change (
    change_by, 
    andon_event_id, 
    cancelled_by,
    cancelled_at,
    status
) 
VALUES (
    $1, $2, $1, NOW(), 'Cancelled'
)
`

	updateQuery := `
UPDATE 
  andon_event
SET 
  cancelled_by = $1,
  cancelled_at = NOW(),
  status = 'Cancelled',
  last_updated = NOW()
WHERE 
  andon_event_id = $2
`

	_, err := exec.Exec(ctx, insertQuery, userID, andonEventID)
	if err != nil {
		return err
	}

	_, err = exec.Exec(ctx, updateQuery, userID, andonEventID)
	if err != nil {
		return err
	}

	return nil
}

func (r *AndonRepository) ReopenAndonEvent(
	ctx context.Context,
	exec db.PGExecutor,
	andonEventID int,
	userID int,
) error {

	insertQuery := `
INSERT INTO 
  andon_change (
    change_by, 
    andon_event_id, 
    raised_by,
    raised_at,
    status
) 
VALUES (
    $1, $2, $1, NOW(), 'Outstanding'
)
`

	updateQuery := `
UPDATE 
  andon_event
SET 
	raised_by = $1,
	raised_at = NOW(),
	acknowledged_by = NULL,
	acknowledged_at = NULL,
	resolved_by = NULL,
	resolved_at = NULL,
	status = 'Outstanding',
	last_updated = NOW()
WHERE 
  andon_event_id = $2
`

	_, err := exec.Exec(ctx, insertQuery, userID, andonEventID)
	if err != nil {
		return err
	}

	_, err = exec.Exec(ctx, updateQuery, userID, andonEventID)
	if err != nil {
		return err
	}

	return nil
}

func generateWhereClause(filters model.ListAndonQuery) (string, []any) {
	var whereClauses []string
	var args []any
	argID := 1

	addInClause := func(column string, values []string) {
		if len(values) == 0 {
			return
		}
		placeholders := make([]string, len(values))
		for i, val := range values {
			args = append(args, val)
			placeholders[i] = fmt.Sprintf("$%d", argID)
			argID++
		}
		whereClauses = append(whereClauses, fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ", ")))
	}

	if filters.StartDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("raised_at >= $%d", argID))
		args = append(args, *filters.StartDate)
		argID++
	}

	if filters.EndDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("raised_at <= $%d", argID))
		args = append(args, *filters.EndDate)
		argID++
	}

	addInClause("issue_name", filters.Issues)
	addInClause("severity", filters.Serverities)
	addInClause("assigned_team_name", filters.Teams)
	addInClause("location", filters.Locations)
	addInClause("status", filters.Statuses)
	addInClause("raised_by_username", filters.RaisedByUsername)
	addInClause("acknowledged_by_username", filters.AcknowledgedByUsername)
	addInClause("resolved_by_username", filters.ResolvedByUsername)

	if len(whereClauses) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(whereClauses, " AND "), args
}
