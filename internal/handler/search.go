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
		// Add more entities as needed
	}

	err := appurl.Unmarshal(r.URL.Query(), &params)
	if err != nil {
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	var allowedSearchEntities []model.SearchEntity
	for _, entity := range searchEntities {
		if entity.HasPermission(ctx.User.Permissions) {
			allowedSearchEntities = append(allowedSearchEntities, entity)
		}
	}

	params.Q = strings.TrimSpace(params.Q)

	results, err := h.searchService.Search(r.Context(), params.Q, allowedSearchEntities, ctx.User.UserID)
	if err != nil {
		_ = searchview.SearchPage(searchview.SearchPageProps{
			Ctx:             ctx,
			SearchTerm:      params.Q,
			SearchEntities:  allowedSearchEntities,
			UserPermissions: ctx.User.Permissions,
		}).
			Render(w)
		return
	}

	_ = searchview.SearchPage(searchview.SearchPageProps{
		Ctx:             ctx,
		SearchTerm:      params.Q,
		SearchEntities:  allowedSearchEntities,
		Results:         results,
		UserPermissions: ctx.User.Permissions,
	}).
		Render(w)

	return
}
