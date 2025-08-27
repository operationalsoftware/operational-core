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
	commentService service.CommentService,
	teamService service.TeamService,
	fileService service.FileService,
) {
	andonHandler := handler.NewAndonHandler(andonService, andonIssueService, commentService, teamService, fileService)

	mux.HandleFunc("GET /andons", andonHandler.HomePage)
	mux.HandleFunc("GET /andons/all", andonHandler.AllAndonsPage)

	mux.HandleFunc("GET /andons/add", andonHandler.AddPage)
	mux.HandleFunc("POST /andons/add", andonHandler.Add)

	mux.HandleFunc("GET /andons/{andonID}", andonHandler.AndonDetailsPage)

	mux.HandleFunc("POST /andons/{entityId}/comments", andonHandler.AddComment)
	mux.HandleFunc("POST /andons/{entityId}/comments/{commentId}/attachment", andonHandler.AddAttachment)

	mux.HandleFunc("POST /andons/{andonID}/{action}/update", andonHandler.AndonUpdate)

}
