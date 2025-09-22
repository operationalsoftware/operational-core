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
	galleryRepository *repository.GalleryRepository
}

func NewAndonService(
	db *pgxpool.Pool,
	swiftConn *swift.Connection,
	andonRepo *repository.AndonRepository,
	commentRepository *repository.CommentRepository,
	galleryRepository *repository.GalleryRepository,
) *AndonService {
	return &AndonService{
		db:                db,
		swiftConn:         swiftConn,
		andonRepository:   andonRepo,
		commentRepository: commentRepository,
		galleryRepository: galleryRepository,
	}
}

func (s *AndonService) CreateAndonEvent(
	ctx context.Context,
	andon model.NewAndon,
	userID int,
) error {
	var err error

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)


	galleryId, err:= s.galleryRepository.CreateGallery(
		ctx,
		tx,
		userID,
	)
	if err != nil {
		return err
	}

	andon.GalleryID = galleryId

	err = s.andonRepository.CreateAndonEvent(
		ctx,
		tx,
		andon,
		userID,
	)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *AndonService) GetAndonByID(
	ctx context.Context,
	andonEventID int,
	userID int,
) (*model.Andon, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	andonEvent, err := s.andonRepository.GetAndonByID(
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

func (s *AndonService) UpdateAndon(
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
		err = s.andonRepository.AcknowledgeAndon(
			ctx,
			tx,
			andonEventID,
			userID,
		)
	case "resolve":
		err = s.andonRepository.ResolveAndon(
			ctx,
			tx,
			andonEventID,
			userID,
		)
	case "cancel":
		err = s.andonRepository.CancelAndon(
			ctx,
			tx,
			andonEventID,
			userID,
		)
	case "reopen":
		err = s.andonRepository.ReopenAndon(
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
) ([]model.Andon, int, model.AndonAvailableFilters, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.Andon{}, 0, model.AndonAvailableFilters{}, err
	}
	defer tx.Rollback(ctx)

	var andons []model.Andon

	andons, err = s.andonRepository.ListAndons(
		ctx,
		tx,
		q,
		userID,
	)
	if err != nil {
		return []model.Andon{}, 0, model.AndonAvailableFilters{}, err
	}

	count, err := s.andonRepository.Count(ctx, tx, q)
	if err != nil {
		return []model.Andon{}, 0, model.AndonAvailableFilters{}, err
	}

	availableFilters, err := s.andonRepository.GetAvailableFilters(ctx, tx, model.AndonFilters{
		StartDate:              q.StartDate,
		EndDate:                q.EndDate,
		LocationIn:              q.LocationIn,
		IssueIn:                 q.IssueIn,
		TeamIn:                  q.TeamIn,
		SeverityIn: q.SeverityIn,
		StatusIn: q.StatusIn,
		RaisedByUsernameIn:       q.RaisedByUsernameIn,
		AcknowledgedByUsernameIn: q.AcknowledgedByUsernameIn,
		ResolvedByUsernameIn:     q.ResolvedByUsernameIn,
	})
	if err != nil {
		return []model.Andon{}, 0, model.AndonAvailableFilters{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return []model.Andon{}, 0, model.AndonAvailableFilters{}, err
	}

	return andons, count, availableFilters, nil
}

func (s *AndonService) GetAndonChangelog(
	ctx context.Context,
	andonEventID int,
	userID int,
) ([]model.AndonChange, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	changelog, err := s.andonRepository.GetAndonChangelog(
		ctx,
		tx,
		andonEventID,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return changelog, nil
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

	_, err = s.commentRepository.AddComment(
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
