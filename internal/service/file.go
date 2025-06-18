package service

import (
	"app/internal/repository"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FileService struct {
	db            *pgxpool.Pool
	pdfRepository *repository.PDFRepository
}

func NewFileService(
	db *pgxpool.Pool,
	pdfRepository *repository.PDFRepository,
) *FileService {
	return &FileService{
		db:            db,
		pdfRepository: pdfRepository,
	}
}

func (s *FileService) GenerateFile(
	ctx context.Context,
	template string,
	inputParams map[string]interface{},
) error {

	return nil
}
