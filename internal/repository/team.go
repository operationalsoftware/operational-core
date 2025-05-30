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
