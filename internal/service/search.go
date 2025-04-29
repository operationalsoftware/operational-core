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
	searchEntities []model.SearchEntity,
	userID int,
) (model.SearchResults, error) {

	results := model.SearchResults{}

	entityNames := make([]string, 0, len(searchEntities))
	for _, e := range searchEntities {
		entityNames = append(entityNames, e.Name)
	}

	searchEntitiesFilter := make(map[string]bool)
	for _, t := range entityNames {
		searchEntitiesFilter[t] = true
	}

	searchInput := model.SearchInput{
		Q: searchTerm,
		E: entityNames,
	}

	if err := s.SearchRepo.CreateRecentSearch(ctx, s.db, searchInput, userID); err != nil {
		log.Printf("Failed to save recent search: %v", err)
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

	return results, nil
}
