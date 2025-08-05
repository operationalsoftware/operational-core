package repository

import (
	"app/internal/model"
	"app/pkg/db"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type TeamRepository struct{}

func NewTeamRepository() *TeamRepository {
	return &TeamRepository{}
}

// Create
func (r *TeamRepository) Create(
	ctx context.Context,
	exec db.PGExecutor,
	team model.NewTeam,
) error {
	query := `
INSERT INTO team (
	team_name
) VALUES ($1)
`
	_, err := exec.Exec(ctx, query, team.TeamName)
	return err
}

// Read
var teamSelect = `
SELECT
	team_id,
	team_name,
	is_archived
`

func (r *TeamRepository) GetByID(
	ctx context.Context,
	exec db.PGExecutor,
	teamID int,
) (*model.Team, error) {
	query := `
` + teamSelect + `

FROM
	team

WHERE
	team_id = $1
`
	var team model.Team
	err := exec.QueryRow(ctx, query, teamID).Scan(
		&team.TeamID,
		&team.TeamName,
		&team.IsArchived,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &team, nil
}

func (r *TeamRepository) List(
	ctx context.Context,
	exec db.PGExecutor,
	q model.ListTeamsQuery,
) ([]model.Team, error) {

	whereClause := r.generateWhereClause(q)
	limit := q.PageSize
	offset := (q.Page - 1) * q.PageSize
	orderByClause, err := q.Sort.ToOrderByClause(model.Team{})

	if orderByClause == "" {
		orderByClause = "ORDER BY team_name ASC"
	}

	query := `
` + teamSelect + `

FROM
	team

` + whereClause + `

` + orderByClause + `

LIMIT $1 OFFSET $2
`
	rows, err := exec.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []model.Team
	for rows.Next() {
		var team model.Team
		if err := rows.Scan(
			&team.TeamID,
			&team.TeamName,
			&team.IsArchived,
		); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func (r *TeamRepository) Count(
	ctx context.Context,
	exec db.PGExecutor,
	q model.ListTeamsQuery,
) (int, error) {

	whereClause := r.generateWhereClause(q)
	query := `
SELECT
	COUNT(*)
FROM
	team
` + whereClause

	var count int
	err := exec.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *TeamRepository) generateWhereClause(q model.ListTeamsQuery) string {
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
func (r *TeamRepository) UpdateTeam(
	ctx context.Context,
	exec db.PGExecutor,
	teamID int,
	update model.TeamUpdate,
) error {
	// Ensure team exists
	existing, err := r.GetByID(ctx, exec, teamID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("team with ID %d not found", teamID)
	}

	query := `
UPDATE team
SET
	team_name = $1,
    is_archived = $2
WHERE
	team_id = $3
`
	_, err = exec.Exec(
		ctx,
		query,

		update.TeamName,
		update.IsArchived,
		teamID,
	)
	return err
}

func (r *TeamRepository) GetTeamUsers(
	ctx context.Context,
	exec db.PGExecutor,
	teamID int,
	q model.ListTeamUsersQuery,
) ([]model.TeamUser, error) {

	limit := q.PageSize
	offset := (q.Page - 1) * q.PageSize
	orderByClause, err := q.Sort.ToOrderByClause(model.TeamUser{})

	if orderByClause == "" {
		orderByClause = "ORDER BY ut.created_at DESC"
	}

	query := `
SELECT
	ut.team_id,
	ut.user_id,
	au.username,
	ut.role
FROM user_team ut
INNER JOIN app_user au ON ut.user_id = au.user_id
WHERE ut.team_id = $1
` + orderByClause + `

LIMIT $2 OFFSET $3
`

	rows, err := exec.Query(ctx, query, teamID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teamUsers []model.TeamUser
	for rows.Next() {
		var u model.TeamUser
		if err := rows.Scan(&u.TeamID, &u.UserID, &u.Username, &u.Role); err != nil {
			return nil, err
		}
		teamUsers = append(teamUsers, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return teamUsers, nil
}

func (r *TeamRepository) GetTeamUsersCount(
	ctx context.Context,
	exec db.PGExecutor,
	teamID int,
) (int, error) {

	query := `
SELECT
	COUNT(*)
FROM user_team
WHERE
	team_id = $1
`
	var count int
	err := exec.QueryRow(ctx, query, teamID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *TeamRepository) AssignUser(
	ctx context.Context,
	exec db.PGExecutor,
	update model.TeamUser,
) error {

	query := `
INSERT INTO user_team (
	team_id,
	user_id,
	role
)
VALUES (
	$1,
	$2,
	$3
)
ON CONFLICT (user_id, team_id)
DO UPDATE SET role = EXCLUDED.role
WHERE user_team.role IS DISTINCT FROM EXCLUDED.role;
`
	_, err := exec.Exec(
		ctx,
		query,

		update.TeamID,
		update.UserID,
		update.Role,
	)
	return err
}

func (r *TeamRepository) DeleteTeamUser(
	ctx context.Context,
	exec db.PGExecutor,
	teamID int,
	userID int,
) error {
	query := `
DELETE FROM
	user_team
WHERE
	team_id = $1 AND user_id = $2
`
	_, err := exec.Exec(ctx, query, teamID, userID)
	return err
}
