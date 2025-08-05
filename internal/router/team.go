package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addTeamRoutes(
	mux *http.ServeMux,
	teamService service.TeamService,
	userService service.UserService,
) {
	teamHandler := handler.NewTeamHandler(teamService, userService)

	mux.HandleFunc("GET /teams", teamHandler.TeamsHomePage)

	mux.HandleFunc("GET /teams/add", teamHandler.AddTeamPage)
	mux.HandleFunc("POST /teams/add", teamHandler.AddTeam)

	mux.HandleFunc("GET /teams/{id}", teamHandler.TeamPage)

	mux.HandleFunc("GET /teams/{id}/edit", teamHandler.EditTeamPage)
	mux.HandleFunc("POST /teams/{id}/edit", teamHandler.EditTeam)

	mux.HandleFunc("GET /teams/{id}/assign-user", teamHandler.AssignUserToTeamPage)
	mux.HandleFunc("POST /teams/{id}/assign-user", teamHandler.AssignUserToTeam)

	mux.HandleFunc("DELETE /teams/{id}/delete/{userID}", teamHandler.DeleteTeamUser)
}
