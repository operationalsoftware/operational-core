package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamService struct {
	db             *pgxpool.Pool
	teamRepository *repository.TeamRepository
}

func NewTeamService(db *pgxpool.Pool, teamRepo *repository.TeamRepository) *TeamService {
	return &TeamService{
		db:             db,
		teamRepository: teamRepo,
	}
}

// Create
func (s *TeamService) Create(
	ctx context.Context,
	team model.NewTeam,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.teamRepository.Create(ctx, tx, team); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// Read
func (s *TeamService) GetByID(
	ctx context.Context,
	teamID int,
) (*model.Team, error) {
	return s.teamRepository.GetByID(ctx, s.db, teamID)
}

func (s *TeamService) List(
	ctx context.Context,
	q model.ListTeamsQuery,
) ([]model.Team, int, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.Team{}, 0, err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	teams, err := s.teamRepository.List(ctx, tx, q)
	if err != nil {
		return []model.Team{}, 0, err
	}

	count, err := s.teamRepository.Count(ctx, tx, q)
	if err != nil {
		return []model.Team{}, 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return []model.Team{}, 0, err
	}

	return teams, count, nil
}

// Update
func (s *TeamService) Update(
	ctx context.Context,
	teamID int,
	update model.TeamUpdate,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.teamRepository.UpdateTeam(ctx, tx, teamID, update); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
