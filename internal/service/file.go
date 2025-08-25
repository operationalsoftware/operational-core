package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ncw/swift/v2"
)

type FileService struct {
	db             *pgxpool.Pool
	swiftConn      *swift.Connection
	fileRepository *repository.FileRepository
}

func NewFileService(
	db *pgxpool.Pool,
	swiftConn *swift.Connection,
	fileRepository *repository.FileRepository,
) *FileService {
	return &FileService{
		db:             db,
		swiftConn:      swiftConn,
		fileRepository: fileRepository,
	}
}

func (s *FileService) CreateFile(
	ctx context.Context,
	file *model.File,
	userID int,
) (*model.File, string, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, "", err
	}
	defer tx.Rollback(ctx)

	newFile, signedURL, err := s.fileRepository.CreateFile(
		ctx,
		tx,
		s.swiftConn,
		file,
		userID,
	)
	if err != nil {
		return nil, "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, "", err
	}

	return newFile, signedURL, nil
}

func (s *FileService) CompleteFileUpload(
	ctx context.Context,
	fileID string,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = s.fileRepository.CompleteFileUpload(
		ctx,
		tx,
		fileID,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
