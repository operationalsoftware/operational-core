package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"
	"log"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SearchService struct {
	db         *pgxpool.Pool
	UserRepo   *repository.UserRepository
	SearchRepo *repository.SearchRepository
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
	searchEntities []string,
	userID int,
) (model.SearchResults, error) {

	results := model.SearchResults{}

	searchEntitiesFilter := make(map[string]bool)
	for _, t := range searchEntities {
		searchEntitiesFilter[t] = true
	}

	searchInput := model.SearchInput{
		Q: searchTerm,
		E: searchEntities,
	}

	recentSearches, err := s.SearchRepo.FetchRecentSearches(ctx, s.db, searchInput, userID)
	if err != nil {
		return results, err
	}
	results.RecentSearches = recentSearches

	if searchTerm == "" {
		return results, err
	}

	// User Search
	if len(searchEntities) == 0 || searchEntitiesFilter["user"] {
		users, err := s.UserRepo.SearchUsers(ctx, s.db, searchTerm)
		if err != nil {
			return results, err
		}

		results.Users = users

		if len(results.Users) > 0 {
			sort.Slice(results.Users, func(i, j int) bool {
				return results.Users[i].Relevance > results.Users[j].Relevance
			})
		}
	}

	// Batch Search
	if len(searchEntities) == 0 || searchEntitiesFilter["batch"] {
		batches, err := s.UserRepo.SearchBatches(ctx, s.db, searchTerm)
		if err != nil {
			return results, err
		}
		results.Batches = batches

		if len(results.Batches) > 0 {
			sort.Slice(results.Batches, func(i, j int) bool {
				return results.Batches[i].Relevance > results.Batches[j].Relevance
			})
		}

	}

	if err := s.SearchRepo.CreateRecentSearch(ctx, s.db, searchInput, userID); err != nil {
		log.Printf("Failed to save recent search: %v", err)
	}

	return results, nil
}
