package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addAndonIssueRoutes(
	mux *http.ServeMux,
	andonIssueService service.AndonIssueService,
	teamService service.TeamService,
) {
	andonIssueHandler := handler.NewAndonIssueHandler(andonIssueService, teamService)

	mux.HandleFunc("GET /andon-issues", andonIssueHandler.HomePage)

	mux.HandleFunc("GET /andon-issues/add", andonIssueHandler.AddPage)
	mux.HandleFunc("POST /andon-issues/add", andonIssueHandler.Add)

	mux.HandleFunc("GET /andon-issues/add-group", andonIssueHandler.AddGroupPage)
	mux.HandleFunc("POST /andon-issues/add-group", andonIssueHandler.AddGroup)

	mux.HandleFunc("GET /andon-issues/{id}", andonIssueHandler.AndonIssuePage)

	mux.HandleFunc("GET /andon-issues/{id}/edit", andonIssueHandler.EditPage)
	mux.HandleFunc("POST /andon-issues/{id}/edit", andonIssueHandler.Edit)

	mux.HandleFunc("GET /andon-issues/group/{id}/edit", andonIssueHandler.EditGroupPage)
	mux.HandleFunc("POST /andon-issues/group/{id}/edit", andonIssueHandler.EditGroup)
}
