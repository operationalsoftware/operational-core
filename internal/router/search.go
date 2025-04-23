package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addSearchRoutes(
	mux *http.ServeMux,
	searchService service.SearchService,
) {
	searchHandler := handler.NewSearchHandler(searchService)

	mux.HandleFunc("GET /search", searchHandler.SearchPage)
	// mux.HandleFunc("GET /search-results", searchHandler.Search)

}
