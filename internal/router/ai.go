package router

import (
	"app/internal/handler"
	"app/pkg/env"
	"net/http"
)

func addAIRoutes(mux *http.ServeMux) {
	if env.IsProduction() {
		return
	}

	aiDocsHandler := handler.NewAIDocsHandler()

	mux.HandleFunc("GET /ai/docs", aiDocsHandler.DocsPage)
	mux.HandleFunc("POST /ai/docs/query", aiDocsHandler.Query)
}
