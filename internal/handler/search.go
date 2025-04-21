package handler

import (
	"app/internal/service"
	"app/internal/views/searchview"
	"app/pkg/reqcontext"
	"encoding/json"
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

	_ = searchview.SearchPage(searchview.SearchPageProps{
		Ctx: ctx,
	}).
		Render(w)

	return
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	types := strings.Split(r.URL.Query().Get("types"), ",")

	results, err := h.searchService.Search(r.Context(), searchTerm, types)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
