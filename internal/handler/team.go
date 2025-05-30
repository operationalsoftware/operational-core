package handler

import (
	"app/internal/model"
	"app/internal/service"
	"app/internal/views/teamview"
	"app/pkg/appsort"
	"app/pkg/appurl"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type TeamHandler struct {
	teamService service.TeamService
}

func NewTeamHandler(teamService service.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

func (h *TeamHandler) TeamsHomePage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var uv teamsHomePageUrlVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	uv.normalise()

	sort := appsort.Sort{}
	err = sort.ParseQueryParam(model.Team{}, uv.Sort)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing sort: %v", err), http.StatusBadRequest)
		return
	}

	teams, count, err := h.teamService.List(r.Context(), model.ListTeamsQuery{
		ShowArchived: uv.ShowArchived,
		Sort:         sort,
		Page:         uv.Page,
		PageSize:     uv.PageSize,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "Error listing teams", http.StatusInternalServerError)
		return
	}

	_ = teamview.TeamsHomePage(&teamview.TeamsHomePageProps{
		Ctx:          ctx,
		ShowArchived: uv.ShowArchived,
		Teams:        teams,
		TeamCount:    count,
		Sort:         sort,
		Page:         uv.Page,
		PageSize:     uv.PageSize,
	}).Render(w)
}

type teamsHomePageUrlVals struct {
	ShowArchived bool
	Sort         string
	Page         int
	PageSize     int
}

func (uv *teamsHomePageUrlVals) normalise() {
	if uv.Page == 0 {
		uv.Page = 1
	}
	if uv.PageSize == 0 {
		uv.PageSize = 50
	}
}

func (h *TeamHandler) TeamPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	teamID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid team id", http.StatusBadRequest)
		return
	}

	team, err := h.teamService.GetByID(r.Context(), teamID)
	if err != nil {
		http.Error(w, "Error fetching team", http.StatusInternalServerError)
		return
	}
	if team == nil {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	}

	_ = teamview.TeamPage(&teamview.TeamPageProps{
		Ctx:  ctx,
		Team: *team,
	}).Render(w)
}

func (h *TeamHandler) AddTeamPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_ = teamview.AddTeamPage(&teamview.AddTeamPageProps{
		Ctx: ctx,
	}).Render(w)
}

func (h *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var fd addTeamFormData
	if err := appurl.Unmarshal(r.Form, &fd); err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd.normalise()

	validationErrors := fd.validate()

	if len(validationErrors) > 0 {
		_ = teamview.AddTeamPage(&teamview.AddTeamPageProps{
			Ctx:              ctx,
			Values:           r.Form,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	err := h.teamService.Create(r.Context(), model.NewTeam{
		TeamName: fd.TeamName,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "Error creating team", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/teams", http.StatusSeeOther)
}

type addTeamFormData struct {
	TeamName string
}

func (fd *addTeamFormData) normalise() {
	fd.TeamName = strings.TrimSpace(fd.TeamName)
}

func (fd *addTeamFormData) validate() validate.ValidationErrors {
	var ve validate.ValidationErrors = make(map[string][]string)

	validate.MinLength(&ve, "TeamName", fd.TeamName, 3)
	validate.MaxLength(&ve, "TeamName", fd.TeamName, 50)

	return ve
}

func (h *TeamHandler) EditTeamPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	teamID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid team id", http.StatusBadRequest)
		return
	}

	team, err := h.teamService.GetByID(r.Context(), teamID)
	if err != nil {
		http.Error(w, "Error getting team", http.StatusInternalServerError)
		return
	}
	if team == nil {
		http.Error(w, "Team does not exist", http.StatusBadRequest)
		return
	}

	_ = teamview.EditTeamPage(&teamview.EditTeamPageProps{
		Ctx:  ctx,
		Team: *team,
	}).Render(w)
}

func (h *TeamHandler) EditTeam(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	teamID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid team id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var fd editTeamFormData
	if err := appurl.Unmarshal(r.Form, &fd); err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd.normalise()

	validationErrors := fd.validate()

	if len(validationErrors) > 0 {
		_ = teamview.EditTeamPage(&teamview.EditTeamPageProps{
			Ctx:              ctx,
			Values:           r.Form,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	err = h.teamService.Update(r.Context(), teamID, model.TeamUpdate{
		TeamName:   fd.TeamName,
		IsArchived: fd.IsArchived,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "Error updating team", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/teams", http.StatusSeeOther)
}

type editTeamFormData struct {
	TeamName   string
	IsArchived bool
}

func (fd *editTeamFormData) normalise() {
	fd.TeamName = strings.TrimSpace(fd.TeamName)
}

func (fd *editTeamFormData) validate() validate.ValidationErrors {
	var ve validate.ValidationErrors = make(map[string][]string)

	validate.MinLength(&ve, "TeamName", fd.TeamName, 3)
	validate.MaxLength(&ve, "TeamName", fd.TeamName, 50)

	return ve
}
