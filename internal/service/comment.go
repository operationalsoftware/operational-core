package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentService struct {
	db                *pgxpool.Pool
	commentRepository *repository.CommentRepository
	fileRepository    *repository.FileRepository
}

func NewCommentService(
	db *pgxpool.Pool,
	commentRepository *repository.CommentRepository,
) *CommentService {
	return &CommentService{
		db:                db,
		commentRepository: commentRepository,
	}
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
