package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addTeamRoutes(
	mux *http.ServeMux,
	teamService service.TeamService,
) {
	teamHandler := handler.NewTeamHandler(teamService)

	mux.HandleFunc("GET /teams", teamHandler.TeamsHomePage)

	mux.HandleFunc("GET /teams/add", teamHandler.AddTeamPage)
	mux.HandleFunc("POST /teams/add", teamHandler.AddTeam)

	mux.HandleFunc("GET /teams/{id}", teamHandler.TeamPage)

	mux.HandleFunc("GET /teams/{id}/edit", teamHandler.EditTeamPage)
	mux.HandleFunc("POST /teams/{id}/edit", teamHandler.EditTeam)
}
