package repository

import (
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
	andon model.NewAndon,
	userID int,
) error {

	andonQuery := `
INSERT INTO andon (
	description,
	andon_issue_id,
	source,
	location,
	raised_by
)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
)
RETURNING andon_id
`

	changeQuery := `
INSERT INTO
  andon_change (
    change_by,
    andon_id,
	description,
    raised_by
)
VALUES ($1, $2, $3, $1)
`

	var newAndonID int
	err := exec.QueryRow(
		ctx, andonQuery,

		andon.Description,
		andon.IssueID,
		andon.Source,
		andon.Location,
		userID,
	).Scan(&newAndonID)
	if err != nil {
		return err
	}

	_, err = exec.Exec(
		ctx, changeQuery,

		userID,
		newAndonID,
		andon.Description,
	)
	if err != nil {
		return err
	}

	return nil
}

func andonSelectClause(currentUserIDPlaceholder int) string {

	currentUserIDPlaceholderStr := fmt.Sprintf("$%d", currentUserIDPlaceholder)

	return `
SELECT
	andon_id,
	description,
	andon_issue_id,
	source,
	location,
	assigned_team,
	assigned_team_name,
	raised_by_username,
	raised_at,
	is_acknowledged,
	acknowledged_by_username,
	acknowledged_at,
	is_resolved,
	resolved_by_username,
	resolved_at,
	is_cancelled,
	cancelled_by_username,
	cancelled_at,
	last_updated,
	name_path,
	severity,
	is_open,
	status,
	(
		is_cancelled = false
	    AND
		is_acknowledged = false
		AND
		(
			(
				severity = 'Self-resolvable'
				AND
				raised_by = ` + currentUserIDPlaceholderStr + `
			)
			OR
			assigned_team IN (SELECT team_id FROM user_team WHERE user_id = ` + currentUserIDPlaceholderStr + `)
		)
	) AS can_user_acknowledge,
	(
		severity <> 'Info'
		AND
		is_cancelled = false
		AND
		is_acknowledged = true
		AND
		is_resolved = false
		AND
		(
			(
				severity = 'Self-resolvable'
				AND
				raised_by = ` + currentUserIDPlaceholderStr + `
			)
			OR
			assigned_team IN (SELECT team_id FROM user_team WHERE user_id = ` + currentUserIDPlaceholderStr + `)
		)
	) AS can_user_resolve,
	(
		is_open = true
		AND
		(
			(
				severity = 'Self-resolvable'
				AND
				raised_by = ` + currentUserIDPlaceholderStr + `
			)
			OR
			assigned_team IN (SELECT team_id FROM user_team WHERE user_id = ` + currentUserIDPlaceholderStr + `)
		)
	) AS can_user_cancel,
	(
		is_cancelled = true
		AND
		is_resolved = true
		AND
		(
			(
				severity = 'Self-resolvable'
				AND
				raised_by = ` + currentUserIDPlaceholderStr + `
			)
			OR
			assigned_team IN (SELECT team_id FROM user_team WHERE user_id = ` + currentUserIDPlaceholderStr + `)
		)
	) AS can_user_reopen
`
}

func (r *AndonRepository) GetAndonByID(
	ctx context.Context,
	exec db.PGExecutor,
	andonEventID int,
	userID int,
) (*model.Andon, error) {

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

	query := andonSelectClause(2) + `
FROM
	andon_view
WHERE
	andon_id = $1
`

	var andon model.Andon
	err = exec.QueryRow(
		ctx, query, andonEventID, userID,
	).Scan(
		&andon.AndonID,
		&andon.Description,
		&andon.AndonIssueID,
		&andon.Source,
		&andon.Location,
		&andon.AssignedTeam,
		&andon.AssignedTeamName,
		&andon.RaisedByUsername,
		&andon.RaisedAt,
		&andon.IsAcknowledged,
		&andon.AcknowledgedByUsername,
		&andon.AcknowledgedAt,
		&andon.IsResolved,
		&andon.ResolvedByUsername,
		&andon.ResolvedAt,
		&andon.IsCancelled,
		&andon.CancelledByUsername,
		&andon.CancelledAt,
		&andon.LastUpdated,
		&andon.NamePath,
		&andon.Severity,
		&andon.IsOpen,
		&andon.Status,
		&andon.CanUserAcknowledge,
		&andon.CanUserResolve,
		&andon.CanUserCancel,
		&andon.CanUserReopen,
	)
	if err != nil {
		return nil, err
	}

	return &andon, err
}

func (r *AndonRepository) ListAndons(
	ctx context.Context,
	exec db.PGExecutor,
	q model.ListAndonQuery,
	currentUserID int,
) ([]model.Andon, error) {

	whereClause, args := generateWhereClause(q)

	currentUserIDPlaceholder := len(args) + 1
	limitPlaceholder := fmt.Sprintf("$%d", len(args)+2)
	offsetPlaceholder := fmt.Sprintf("$%d", len(args)+3)

	query := andonSelectClause(currentUserIDPlaceholder) + `
FROM andon_view
`

	limit := q.PageSize
	offset := (q.Page - 1) * q.PageSize
	orderByClause, _ := q.Sort.ToOrderByClause(model.Andon{})
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

	var andons []model.Andon

	for rows.Next() {
		var andon model.Andon

		if err := rows.Scan(
			&andon.AndonID,
			&andon.Description,
			&andon.AndonIssueID,
			&andon.Source,
			&andon.Location,
			&andon.AssignedTeam,
			&andon.AssignedTeamName,
			&andon.RaisedByUsername,
			&andon.RaisedAt,
			&andon.IsAcknowledged,
			&andon.AcknowledgedByUsername,
			&andon.AcknowledgedAt,
			&andon.IsResolved,
			&andon.ResolvedByUsername,
			&andon.ResolvedAt,
			&andon.IsCancelled,
			&andon.CancelledByUsername,
			&andon.CancelledAt,
			&andon.LastUpdated,
			&andon.NamePath,
			&andon.Severity,
			&andon.IsOpen,
			&andon.Status,
			&andon.CanUserAcknowledge,
			&andon.CanUserResolve,
			&andon.CanUserCancel,
			&andon.CanUserReopen,
		); err != nil {
			return nil, err
		}

		andons = append(andons, andon)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return andons, nil
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
			Severities:             baseFilters.Severities,
			Teams:                  baseFilters.Teams,
			Locations:              baseFilters.Locations,
			RaisedByUsername:       baseFilters.RaisedByUsername,
			AcknowledgedByUsername: baseFilters.AcknowledgedByUsername,
			ResolvedByUsername:     baseFilters.ResolvedByUsername,
		}

		switch key {
		case "IssueIn":
			queryFilters.Issues = nil
		case "SeverityIn":
			queryFilters.Severities = nil
		case "TeamIn":
			queryFilters.Teams = nil
		case "LocationIn":
			queryFilters.Locations = nil
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

func (r *AndonRepository) GetAndonChangelog(
	ctx context.Context,
	exec db.PGExecutor,
	andonID int,
) ([]model.AndonChange, error) {

	query := `
SELECT
	andon_id,
	andon_change_id,
	change_by,
	change_by_username,
	change_at,
    is_creation,
	description,
	raised_by,
	raised_by_username,
	raised_at,
	acknowledged_by,
	acknowledged_by_username,
	acknowledged_at,
	resolved_by,
	resolved_by_username,
	resolved_at,
	cancelled_by,
	cancelled_by_username,
	cancelled_at
FROM
	andon_change_view

WHERE andon_id = $1
ORDER BY change_at DESC
`

	rows, err := exec.Query(ctx, query, andonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var changes []model.AndonChange
	for rows.Next() {
		var change model.AndonChange
		if err := rows.Scan(
			&change.AndonID,
			&change.AndonChangeID,
			&change.ChangeBy,
			&change.ChangeByUsername,
			&change.ChangeAt,
			&change.IsCreation,
			&change.Description,
			&change.RaisedBy,
			&change.RaisedByUsername,
			&change.RaisedAt,
			&change.AcknowledgedBy,
			&change.AcknowledgedByUsername,
			&change.AcknowledgedAt,
			&change.ResolvedBy,
			&change.ResolvedByUsername,
			&change.ResolvedAt,
			&change.CancelledBy,
			&change.CancelledByUsername,
			&change.CancelledAt,
		); err != nil {
			return nil, err
		}
		changes = append(changes, change)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return changes, nil
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
    andon_id, 
    acknowledged_by
) 
VALUES ($1, $2, $1)
`

	updateQuery := `
UPDATE 
  andon
SET 
  acknowledged_by = $1,
  acknowledged_at = NOW(),
  last_updated = NOW()
WHERE 
  andon_id = $2
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
    andon_id, 
    resolved_by, 
) 
VALUES ($1, $2, $1)
`

	updateQuery := `
UPDATE 
  andon
SET 
  resolved_by = $1,
  resolved_at = NOW(),
  last_updated = NOW()
WHERE 
  andon_id = $2
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
    andon_id, 
    cancelled_by,
    cancelled_at,
    status
) 
VALUES ($1, $2, $1)
`

	updateQuery := `
UPDATE 
  andon
SET 
  cancelled_by = $1,
  cancelled_at = NOW(),
  last_updated = NOW()
WHERE 
  andon_id = $2
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
    andon_id, 
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
  andon
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
  andon_id = $2
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

	if filters.IsAcknowledged != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("is_acknowledged = $%d", argID))
		args = append(args, *filters.IsAcknowledged)
		argID++
	}

	if filters.IsResolved != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("is_resolved = $%d", argID))
		args = append(args, *filters.IsResolved)
		argID++
	}

	if filters.IsCancelled != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("is_cancelled = $%d", argID))
		args = append(args, *filters.IsCancelled)
		argID++
	}

	if filters.IsOpen != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("is_open = $%d", argID))
		args = append(args, *filters.IsOpen)
		argID++
	}

	addInClause("issue_name", filters.Issues)
	addInClause("severity", filters.Severities)
	addInClause("assigned_team_name", filters.Teams)
	addInClause("location", filters.Locations)
	addInClause("raised_by_username", filters.RaisedByUsername)
	addInClause("acknowledged_by_username", filters.AcknowledgedByUsername)
	addInClause("resolved_by_username", filters.ResolvedByUsername)

	if len(whereClauses) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(whereClauses, " AND "), args
}
