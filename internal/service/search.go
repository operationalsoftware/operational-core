package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"
	"log"
	"sort"
	"sync"

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

	recentSearches, err := s.SearchRepo.FetchRecentSearches(ctx, s.db, searchInput, userID)
	if err != nil {
		return results, err
	}
	results.RecentSearches = recentSearches

	if searchTerm == "" {
		return results, err
	}

	if err := s.SearchRepo.CreateRecentSearch(ctx, s.db, searchInput, userID); err != nil {
		log.Printf("Failed to save recent search: %v", err)
	}

	// Goroutine variables
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		errChan = make(chan error, 5)
	)

	// User Search
	if len(searchEntities) == 0 || searchEntitiesFilter["user"] {
		wg.Add(1)

		go func() {
			defer wg.Done()

			users, err := s.UserRepo.SearchUsers(ctx, s.db, searchTerm)
			if err != nil {
				errChan <- err
				return
			}

			mu.Lock()
			results.Users = users
			mu.Unlock()

			if len(results.Users) > 0 {
				sort.Slice(results.Users, func(i, j int) bool {
					return results.Users[i].Relevance > results.Users[j].Relevance
				})
			}
		}()
	}

	// Wait for all to complete
	wg.Wait()
	close(errChan)

	// Collect errors, return the first one if any
	var firstErr error
	for err := range errChan {
		if firstErr == nil {
			firstErr = err
		}
	}
	if firstErr != nil {
		return results, firstErr
	}

	return results, nil
}
