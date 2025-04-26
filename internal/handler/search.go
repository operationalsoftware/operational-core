package handler

import (
	"app/internal/model"
	"app/internal/service"
	"app/internal/views/searchview"
	"app/pkg/appurl"
	"app/pkg/reqcontext"
	"net/http"
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

	err := appurl.Unmarshal(r.URL.Query(), &params)
	if err != nil {
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	results, err := h.searchService.Search(r.Context(), params.Q, params.E, ctx.User.UserID)
	if err != nil {
		_ = searchview.SearchPage(searchview.SearchPageProps{
			Ctx:            ctx,
			SearchTerm:     params.Q,
			SearchEntities: params.E,
		}).
			Render(w)
		return
	}

	_ = searchview.SearchPage(searchview.SearchPageProps{
		Ctx:            ctx,
		SearchTerm:     params.Q,
		SearchEntities: params.E,
		Results:        results,
	}).
		Render(w)

	return
}
