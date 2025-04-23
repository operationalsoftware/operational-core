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

	results, err := h.searchService.Search(r.Context(), params.Q, params.E)
	if err != nil {
		_ = searchview.SearchPage(searchview.SearchPageProps{
			Ctx: ctx,
		}).
			Render(w)
	}

	_ = searchview.SearchPage(searchview.SearchPageProps{
		Ctx:     ctx,
		Results: results,
	}).
		Render(w)

	return
}

// func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {

// 	results, err := h.searchService.Search(r.Context(), searchTerm, types)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(results)
// }
