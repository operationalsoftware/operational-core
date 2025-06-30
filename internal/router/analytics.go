package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addAnalyticsRoutes(
	mux *http.ServeMux,
	analyticsService service.AnalyticsService,
) {
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	mux.HandleFunc("GET /analytics/dashboard", analyticsHandler.DashboardHandler)
	mux.HandleFunc("GET /analytics/stats", analyticsHandler.StatsHandler)

}
