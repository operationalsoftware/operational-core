package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamService struct {
	db             *pgxpool.Pool
	teamRepository *repository.TeamRepository
	userRepository *repository.UserRepository
}

func NewTeamService(db *pgxpool.Pool,
	teamRepo *repository.TeamRepository,
	userRepo *repository.UserRepository) *TeamService {
	return &TeamService{
		db:             db,
		teamRepository: teamRepo,
		userRepository: userRepo,
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

func (s *TeamService) GetTeamUsers(
	ctx context.Context,
	teamID int,
	q model.ListTeamUsersQuery,
) ([]model.TeamUser, int, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	teamUsers, err := s.teamRepository.GetTeamUsers(ctx, tx, teamID, q)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.teamRepository.GetTeamUsersCount(ctx, tx, teamID)
	if err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, 0, err
	}

	return teamUsers, count, nil
}

func (s *TeamService) AssignUser(
	ctx context.Context,
	update model.TeamUser,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	existing, err := s.userRepository.GetUserByID(ctx, tx, update.UserID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("user with ID %d not found", update.UserID)
	}

	existingTeam, err := s.teamRepository.GetByID(ctx, tx, update.TeamID)
	if err != nil {
		return err
	}
	if existingTeam == nil {
		return fmt.Errorf("team with ID %d not found", update.UserID)
	}

	if err := s.teamRepository.AssignUser(ctx, tx, update); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *TeamService) DeleteTeamUser(
	ctx context.Context,
	update model.TeamUser,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.teamRepository.DeleteTeamUser(ctx, tx, update.TeamID, update.UserID); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
