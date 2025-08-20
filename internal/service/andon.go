package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ncw/swift/v2"
)

type AndonService struct {
	db                *pgxpool.Pool
	swiftConn         *swift.Connection
	andonRepository   *repository.AndonRepository
	commentRepository *repository.CommentRepository
	fileRepository    *repository.FileRepository
}

func NewAndonService(
	db *pgxpool.Pool,
	swiftConn *swift.Connection,
	andonRepo *repository.AndonRepository,
	commentRepository *repository.CommentRepository,
	fileRepository *repository.FileRepository,
) *AndonService {
	return &AndonService{
		db:                db,
		swiftConn:         swiftConn,
		andonRepository:   andonRepo,
		commentRepository: commentRepository,
		fileRepository:    fileRepository,
	}
}

func (s *AndonService) CreateAndonEvent(
	ctx context.Context,
	andonIssue model.NewAndonEvent,
	userID int,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.andonRepository.CreateAndonEvent(
		ctx,
		s.db,
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

func (s *AndonService) GetAndonEventByID(
	ctx context.Context,
	andonEventID int,
	userID int,
) (*model.AndonEvent, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	andonEvent, err := s.andonRepository.GetAndonEventByID(
		ctx,
		s.db,
		andonEventID,
		userID,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return andonEvent, nil
}

func (s *AndonService) UpdateAndonEvent(
	ctx context.Context,
	andonEventID int,
	action string,
	userID int,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	switch action {
	case "acknowledge":
		err = s.andonRepository.AcknowledgeAndonEvent(
			ctx,
			tx,
			andonEventID,
			userID,
		)
	case "resolve":
		err = s.andonRepository.ResolveAndonEvent(
			ctx,
			tx,
			andonEventID,
			userID,
		)
	case "cancel":
		err = s.andonRepository.CancelAndonEvent(
			ctx,
			tx,
			andonEventID,
			userID,
		)
	case "reopen":
		err = s.andonRepository.ReopenAndonEvent(
			ctx,
			tx,
			andonEventID,
			userID,
		)
	}

	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *AndonService) ListAndons(
	ctx context.Context,
	q model.ListAndonQuery,
	userID int,
) ([]model.AndonEvent, int, model.AndonAvailableFilters, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.AndonEvent{}, 0, model.AndonAvailableFilters{}, err
	}
	defer tx.Rollback(ctx)

	var andonEvents []model.AndonEvent

	andonEvents, err = s.andonRepository.ListAndons(
		ctx,
		tx,
		q,
		userID,
	)
	if err != nil {
		return []model.AndonEvent{}, 0, model.AndonAvailableFilters{}, err
	}

	count, err := s.andonRepository.Count(ctx, tx, q)
	if err != nil {
		return []model.AndonEvent{}, 0, model.AndonAvailableFilters{}, err
	}

	filters, err := s.andonRepository.GetAvailableFilters(ctx, tx, model.AndonFilters{
		StartDate:              q.StartDate,
		EndDate:                q.EndDate,
		Issues:                 q.Issues,
		Teams:                  q.Teams,
		Locations:              q.Locations,
		Statuses:               q.Statuses,
		RaisedByUsername:       q.RaisedByUsername,
		AcknowledgedByUsername: q.AcknowledgedByUsername,
		ResolvedByUsername:     q.ResolvedByUsername,
	})
	if err != nil {
		return []model.AndonEvent{}, 0, model.AndonAvailableFilters{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return []model.AndonEvent{}, 0, model.AndonAvailableFilters{}, err
	}

	return andonEvents, count, filters, nil
}

func (s *AndonService) GetAndonByID(
	ctx context.Context,
	andonEventID int,
	userID int,
) ([]model.AndonChange, []model.Comment, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(ctx)

	andonChanges, err := s.andonRepository.GetAndonChangeLog(
		ctx,
		tx,
		andonEventID,
	)
	if err != nil {
		return nil, nil, err
	}

	andonComments, err := s.commentRepository.GetComments(
		ctx,
		tx,
		"andon",
		andonEventID,
	)
	if err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, err
	}

	return andonChanges, andonComments, nil
}

func (s *AndonService) CreateAndonComment(
	ctx context.Context,
	comment *model.NewComment,
	userID int,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = s.commentRepository.AddComment(
		ctx,
		tx,
		comment,
		userID,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
