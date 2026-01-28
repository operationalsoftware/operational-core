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
	SearchRepo *repository.SearchRepository
}

func NewSearchService(
	db *pgxpool.Pool,
	searchRepository *repository.SearchRepository,
) *SearchService {
	return &SearchService{
		db:         db,
		SearchRepo: searchRepository,
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
	group := newSearchRoutineGroup()
	var mu sync.Mutex

	// User Search
	if len(searchEntities) == 0 || searchEntitiesFilter["user"] {
		group.run(func() error {
			users, err := s.SearchRepo.SearchUsers(ctx, s.db, searchTerm)
			if err != nil {
				return err
			}

			mu.Lock()
			results.Users = users
			mu.Unlock()

			if len(results.Users) > 0 {
				sort.Slice(results.Users, func(i, j int) bool {
					return results.Users[i].Relevance > results.Users[j].Relevance
				})
			}

			return nil
		})
	}

	// Stock Item Search
	if len(searchEntities) == 0 || searchEntitiesFilter["stock-item"] {
		group.run(func() error {
			stockItems, err := s.SearchRepo.SearchStockItems(ctx, s.db, searchTerm)
			if err != nil {
				return err
			}

			mu.Lock()
			results.StockItems = stockItems
			mu.Unlock()

			if len(results.StockItems) > 0 {
				sort.Slice(results.StockItems, func(i, j int) bool {
					return results.StockItems[i].Relevance > results.StockItems[j].Relevance
				})
			}

			return nil
		})
	}

	// Resource Search
	if len(searchEntities) == 0 || searchEntitiesFilter["resource"] {
		group.run(func() error {
			resources, err := s.SearchRepo.SearchResources(ctx, s.db, searchTerm)
			if err != nil {
				return err
			}

			mu.Lock()
			results.Resources = resources
			mu.Unlock()

			if len(results.Resources) > 0 {
				sort.Slice(results.Resources, func(i, j int) bool {
					return results.Resources[i].Relevance > results.Resources[j].Relevance
				})
			}

			return nil
		})
	}

	// Service Search
	if len(searchEntities) == 0 || searchEntitiesFilter["service"] {
		group.run(func() error {
			services, err := s.SearchRepo.SearchServices(ctx, s.db, searchTerm)
			if err != nil {
				return err
			}

			mu.Lock()
			results.Services = services
			mu.Unlock()

			if len(results.Services) > 0 {
				sort.Slice(results.Services, func(i, j int) bool {
					return results.Services[i].Relevance > results.Services[j].Relevance
				})
			}

			return nil
		})
	}

	if err := group.wait(); err != nil {
		return results, err
	}

	return results, nil
}

type searchRoutineGroup struct {
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
}

func newSearchRoutineGroup() *searchRoutineGroup {
	return &searchRoutineGroup{}
}

func (g *searchRoutineGroup) run(routine func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := routine(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
			})
		}
	}()
}

func (g *searchRoutineGroup) wait() error {
	g.wg.Wait()
	return g.err
}
