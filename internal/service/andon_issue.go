package service

import (
	"app/internal/model"
	"app/internal/repository"
	"app/pkg/validate"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AndonIssueService struct {
	db                   *pgxpool.Pool
	andonIssueRepository *repository.AndonIssueRepository
}

func NewAndonIssueService(
	db *pgxpool.Pool,
	andonIssueRepo *repository.AndonIssueRepository,
) *AndonIssueService {
	return &AndonIssueService{
		db:                   db,
		andonIssueRepository: andonIssueRepo,
	}
}

// Create
func (s *AndonIssueService) Create(
	ctx context.Context,
	andonIssue model.NewAndonIssue,
	userID int,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.andonIssueRepository.Create(
		ctx,
		tx,
		andonIssue,
		userID,
	); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// Create Group
func (s *AndonIssueService) CreateGroup(
	ctx context.Context,
	andonIssueGroup model.NewAndonIssueGroup,
	userID int,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.andonIssueRepository.CreateGroup(
		ctx,
		tx,
		andonIssueGroup,
		userID,
	); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// Read
func (s *AndonIssueService) GetByID(
	ctx context.Context,
	andonIssueID int,
) (*model.AndonIssue, error) {
	return s.andonIssueRepository.GetByID(ctx, s.db, andonIssueID)
}

func (s *AndonIssueService) List(
	ctx context.Context,
	q model.ListAndonIssuesQuery,
) ([]model.AndonIssue, int, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.AndonIssue{}, 0, err
	}
	defer tx.Rollback(ctx)

	andonIssues, err := s.andonIssueRepository.List(ctx, tx, q)
	if err != nil {
		return []model.AndonIssue{}, 0, err
	}

	count, err := s.andonIssueRepository.Count(ctx, tx, q)
	if err != nil {
		return []model.AndonIssue{}, 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return []model.AndonIssue{}, 0, err
	}

	return andonIssues, count, nil
}

func (s *AndonIssueService) ListGroups(
	ctx context.Context,
) ([]model.AndonIssueGroup, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.AndonIssueGroup{}, err
	}
	defer tx.Rollback(ctx)

	andonGroups, err := s.andonIssueRepository.ListGroups(ctx, tx)
	if err != nil {
		return []model.AndonIssueGroup{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return []model.AndonIssueGroup{}, err
	}

	return andonGroups, nil
}

func (s *AndonIssueService) ListTopLevelGroups(
	ctx context.Context,
) ([]model.AndonIssueGroup, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.AndonIssueGroup{}, err
	}
	defer tx.Rollback(ctx)

	andonGroups, err := s.andonIssueRepository.ListTopLevelGroups(ctx, tx)
	if err != nil {
		return []model.AndonIssueGroup{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return []model.AndonIssueGroup{}, err
	}

	return andonGroups, nil
}

// Update
func (s *AndonIssueService) Update(
	ctx context.Context,
	andonIssueID int,
	update model.AndonIssueUpdate,
	userID int,
) (*validate.ValidationErrors, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	existing, _ := s.andonIssueRepository.GetByID(ctx, tx, andonIssueID)

	if existing == nil {
		return nil, fmt.Errorf("andon issue with ID %d doesn't exist", andonIssueID)
	}

	// validate we aren't archiving a team with active children
	if !existing.IsArchived && update.IsArchived {

		hasActiveChild, err := s.andonIssueRepository.HasActiveChildIssues(ctx, tx, andonIssueID)
		if err != nil {
			return nil, fmt.Errorf("error checking if andon issue has active child issues: %v", err)
		}
		if hasActiveChild {
			validationErrors := validate.ValidationErrors{}
			validationErrors.Add("IsArchived", "Cannot archive an andon issue which has active children")

			return &validationErrors, nil
		}
	}

	if err := s.andonIssueRepository.Update(
		ctx,
		tx,
		andonIssueID,
		update,
		userID,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *AndonIssueService) ListIssuesAndGroups(
	ctx context.Context,
	q model.ListAndonIssuesQuery,
) ([]model.AndonIssueNode, int, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.AndonIssueNode{}, 0, err
	}
	defer tx.Rollback(ctx)

	andonIssues, err := s.andonIssueRepository.ListIssuesAndGroups(ctx, tx, q)
	if err != nil {
		return []model.AndonIssueNode{}, 0, err
	}

	count, err := s.andonIssueRepository.Count(ctx, tx, q)
	if err != nil {
		return []model.AndonIssueNode{}, 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return []model.AndonIssueNode{}, 0, err
	}

	return andonIssues, count, nil
}

func (s *AndonIssueService) GetIssueHierarchy(
	ctx context.Context,
	issueID int,
) ([]int, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	andonIssues, err := s.andonIssueRepository.GetIssueHierarchy(ctx, tx, issueID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return andonIssues, nil
}
