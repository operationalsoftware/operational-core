package handler

import (
	"app/internal/service"
	"app/internal/views/searchview"
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

	_ = searchview.SearchPage(searchview.SearchPageProps{
		Ctx: ctx,
	}).
		Render(w)

	return
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	// ctx := reqcontext.GetContext(r)

	// r.URL.
}
