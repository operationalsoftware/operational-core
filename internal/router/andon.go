package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addAndonRoutes(
	mux *http.ServeMux,
	andonService service.AndonService,
	andonIssueService service.AndonIssueService,
	teamService service.TeamService,
) {
	andonHandler := handler.NewAndonHandler(andonService, andonIssueService, teamService)

	mux.HandleFunc("GET /andons", andonHandler.HomePage)
	mux.HandleFunc("GET /andons/all", andonHandler.AllAndonsPage)

	mux.HandleFunc("GET /andons/add", andonHandler.AddPage)
	mux.HandleFunc("POST /andons/add", andonHandler.Add)

	mux.HandleFunc("GET /andons/{andonID}", andonHandler.AndonDetailsPage)
	mux.HandleFunc("POST /andons/add/comment", andonHandler.AddComment)

	mux.HandleFunc("POST /andons/update/{andonID}/{action}", andonHandler.AndonUpdate)

}
