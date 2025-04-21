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
	authHandler := handler.NewSearchHandler(searchService)

	mux.HandleFunc("GET /search", authHandler.SearchPage)
	mux.HandleFunc("GET /search-results", authHandler.Search)

}
