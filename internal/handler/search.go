package handler

import (
	"app/internal/model"
	"app/internal/service"
	"app/internal/views/searchview"
	"app/pkg/appurl"
	"app/pkg/reqcontext"
	"net/http"
	"strings"
)

type SearchHandler struct {
	searchService service.SearchService
}

func NewSearchHandler(searchService service.SearchService) *SearchHandler {
	return &SearchHandler{searchService: searchService}
}

func (h *SearchHandler) SearchPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	var params model.SearchInput

	var searchEntities = []model.SearchEntity{
		{
			Name:  "user",
			Title: "Users",
			HasPermission: func(permissions model.UserPermissions) bool {
				return permissions.UserAdmin.Access
			},
		},
		{
			Name:  "stock-item",
			Title: "Stock Items",
			HasPermission: func(permissions model.UserPermissions) bool {
				return true
			},
		},
		// Add more entities as needed
	}

	err := appurl.Unmarshal(r.URL.Query(), &params)
	if err != nil {
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	allowedSearchEntities, allowedSearchEntityMap := filterSearchEntities(searchEntities, ctx.User.Permissions)

	params.Q = strings.TrimSpace(params.Q)

	selectedSearchEntities, selectedEntityNames := selectedSearchEntities(
		params.E,
		allowedSearchEntities,
		allowedSearchEntityMap,
	)

	results, err := h.searchService.Search(r.Context(), params.Q, selectedSearchEntities, ctx.User.UserID)
	if err != nil {
		_ = searchview.SearchPage(searchview.SearchPageProps{
			Ctx:              ctx,
			SearchTerm:       params.Q,
			SearchEntities:   allowedSearchEntities,
			SelectedEntities: selectedEntityNames,
			UserPermissions:  ctx.User.Permissions,
		}).
			Render(w)
		return
	}

	_ = searchview.SearchPage(searchview.SearchPageProps{
		Ctx:              ctx,
		SearchTerm:       params.Q,
		SearchEntities:   allowedSearchEntities,
		SelectedEntities: selectedEntityNames,
		Results:          results,
		UserPermissions:  ctx.User.Permissions,
	}).
		Render(w)
}

func filterSearchEntities(
	entities []model.SearchEntity,
	permissions model.UserPermissions,
) ([]model.SearchEntity, map[string]model.SearchEntity) {
	var allowed []model.SearchEntity
	allowedMap := make(map[string]model.SearchEntity, len(entities))
	for _, entity := range entities {
		if entity.HasPermission != nil && !entity.HasPermission(permissions) {
			continue
		}
		allowed = append(allowed, entity)
		allowedMap[entity.Name] = entity
	}

	return allowed, allowedMap
}

func selectedSearchEntities(
	rawEntities []string,
	allowedEntities []model.SearchEntity,
	allowedEntityMap map[string]model.SearchEntity,
) ([]model.SearchEntity, []string) {
	if len(rawEntities) == 0 {
		entityNames := make([]string, 0, len(allowedEntities))
		for _, entity := range allowedEntities {
			entityNames = append(entityNames, entity.Name)
		}
		return allowedEntities, entityNames
	}

	var selected []model.SearchEntity
	var selectedNames []string
	seen := make(map[string]struct{})
	for _, entityValue := range rawEntities {
		for _, entityName := range strings.Split(entityValue, ",") {
			entityName = strings.TrimSpace(entityName)
			if entityName == "" {
				continue
			}

			entity, ok := allowedEntityMap[entityName]
			if !ok {
				continue
			}
			if _, ok := seen[entityName]; ok {
				continue
			}
			seen[entityName] = struct{}{}
			selected = append(selected, entity)
			selectedNames = append(selectedNames, entityName)
		}
	}

	return selected, selectedNames
}
