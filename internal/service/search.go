package service

import (
	"app/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SearchService struct {
	db               *pgxpool.Pool
	searchRepository *repository.AuthRepository
}
