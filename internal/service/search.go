package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SearchService struct {
	db       *pgxpool.Pool
	UserRepo *repository.UserRepository
}

func NewSearchService(
	db *pgxpool.Pool,
	UserRepository *repository.UserRepository,
) *SearchService {
	return &SearchService{
		db:       db,
		UserRepo: UserRepository,
	}
}

func (s *SearchService) Search(
	ctx context.Context,
	searchTerm string,
	types []string,
) (map[string][]model.SearchResult, error) {

	results := make(map[string][]model.SearchResult)

	// You can use a set or map for type filtering
	typeFilter := make(map[string]bool)
	for _, t := range types {
		typeFilter[strings.ToLower(t)] = true
	}

	// User Search
	if len(types) == 0 || typeFilter["user"] {
		users, err := s.UserRepo.SearchUsers(ctx, s.db, searchTerm)
		if err != nil {
			return nil, err
		}
		for _, u := range users {
			results["users"] = append(results["users"], model.SearchResult{
				Type:      "user",
				Data:      u,
				Relevance: u.Relevance,
			})
		}
	}

	// Batch Search
	if len(types) == 0 || typeFilter["batch"] {
		batches, err := s.UserRepo.SearchBatches(ctx, s.db, searchTerm)
		if err != nil {
			return nil, err
		}
		for _, batch := range batches {
			results["batches"] = append(results["batches"], model.SearchResult{
				Type:      "batch",
				Data:      batch,
				Relevance: batch.Relevance,
			})
		}
	}

	// Sort by Relevance descending
	for _, resultList := range results {
		sort.Slice(resultList, func(i, j int) bool {
			return resultList[i].Relevance > resultList[j].Relevance
		})
	}

	return results, nil
}
