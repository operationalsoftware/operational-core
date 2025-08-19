package service

import (
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

func (s *FileService) GenerateFile(
	ctx context.Context,
	template string,
	inputParams map[string]interface{},
) error {

	return nil
}
