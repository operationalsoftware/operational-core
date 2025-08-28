package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ncw/swift/v2"
)

type CommentService struct {
	db                *pgxpool.Pool
	swiftConn         *swift.Connection
	commentRepository *repository.CommentRepository
	fileRepository    *repository.FileRepository
}

func NewCommentService(
	db *pgxpool.Pool,
	swiftConn *swift.Connection,
	commentRepository *repository.CommentRepository,
) *CommentService {
	return &CommentService{
		db:                db,
		swiftConn:         swiftConn,
		commentRepository: commentRepository,
	}
}

func (s *CommentService) GetComments(
	ctx context.Context,
	entity string,
	entityID int,
	userID int,
) ([]model.Comment, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.Comment{}, err
	}
	defer tx.Rollback(ctx)

	comments, err := s.commentRepository.GetComments(
		ctx,
		tx,
		s.swiftConn,
		"Andon",
		entityID,
	)
	if err != nil {
		return []model.Comment{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return []model.Comment{}, err
	}

	return comments, nil
}

func (s *CommentService) CreateComment(
	ctx context.Context,
	comment *model.NewComment,
	userID int,
) (int, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	commentId, err := s.commentRepository.AddComment(
		ctx,
		tx,
		comment,
		userID,
	)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return commentId, nil
}
