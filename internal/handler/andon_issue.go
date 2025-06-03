package handler

import (
	"app/internal/model"
	"app/internal/service"
	"app/internal/views/andonissueview"
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

type AndonIssueHandler struct {
	andonIssueService service.AndonIssueService
	teamService       service.TeamService
}

func NewAndonIssueHandler(
	andonIssueService service.AndonIssueService,
	teamService service.TeamService,
) *AndonIssueHandler {
	return &AndonIssueHandler{
		andonIssueService: andonIssueService,
		teamService:       teamService,
	}
}

func (h *AndonIssueHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var uv andonIssuesHomePageUrlVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	uv.normalise()

	sort := appsort.Sort{}
	err = sort.ParseQueryParam(model.AndonIssue{}, uv.Sort)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing sort: %v", err), http.StatusBadRequest)
		return
	}

	andonIssues, count, err := h.andonIssueService.List(r.Context(), model.ListAndonIssuesQuery{
		ShowArchived: uv.ShowArchived,
		Sort:         sort,
		Page:         uv.Page,
		PageSize:     uv.PageSize,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "Error listing andon issues", http.StatusInternalServerError)
		return
	}

	_ = andonissueview.HomePage(&andonissueview.HomePageProps{
		Ctx:             ctx,
		ShowArchived:    uv.ShowArchived,
		AndonIssues:     andonIssues,
		AndonIssueCount: count,
		Sort:            sort,
		Page:            uv.Page,
		PageSize:        uv.PageSize,
	}).Render(w)
}

type andonIssuesHomePageUrlVals struct {
	ShowArchived bool
	Sort         string
	Page         int
	PageSize     int
}

func (uv *andonIssuesHomePageUrlVals) normalise() {
	if uv.Page == 0 {
		uv.Page = 1
	}
	if uv.PageSize == 0 {
		uv.PageSize = 50
	}
}

func (h *AndonIssueHandler) AndonIssuePage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	andonIssueID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid andon issue ID", http.StatusBadRequest)
		return
	}

	andonIssue, err := h.andonIssueService.GetByID(r.Context(), andonIssueID)
	if err != nil {
		http.Error(w, "Error fetching andon issue", http.StatusInternalServerError)
		return
	}
	if andonIssue == nil {
		http.Error(w, "Andon issue not found", http.StatusNotFound)
		return
	}

	_ = andonissueview.AndonIssuePage(&andonissueview.AndonIssuePageProps{
		Ctx:        ctx,
		AndonIssue: *andonIssue,
	}).Render(w)
}

func (h *AndonIssueHandler) AddPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	andonIssues, _, err := h.andonIssueService.List(
		r.Context(),
		model.ListAndonIssuesQuery{
			Page: 1, PageSize: 10000,
		})
	if err != nil {
		http.Error(w, "Error fetching andon issues", http.StatusInternalServerError)
		return
	}

	teams, _, err := h.teamService.List(r.Context(), model.ListTeamsQuery{
		Page: 1, PageSize: 10000,
	})
	if err != nil {
		http.Error(w, "Error fetching teams", http.StatusInternalServerError)
		return
	}

	_ = andonissueview.AddPage(&andonissueview.AddPageProps{
		Ctx:         ctx,
		AndonIssues: andonIssues,
		Teams:       teams,
	}).Render(w)
}

func (h *AndonIssueHandler) Add(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var fd addAndonIssueFormData
	if err := appurl.Unmarshal(r.Form, &fd); err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd.normalise()

	validationErrors := fd.validate()

	if len(validationErrors) > 0 {
		andonIssues, _, err := h.andonIssueService.List(
			r.Context(),
			model.ListAndonIssuesQuery{
				Page: 1, PageSize: 10000,
			})
		if err != nil {
			http.Error(w, "Error fetching andon issues", http.StatusInternalServerError)
			return
		}

		teams, _, err := h.teamService.List(r.Context(), model.ListTeamsQuery{
			Page: 1, PageSize: 10000,
		})
		if err != nil {
			http.Error(w, "Error fetching teams", http.StatusInternalServerError)
			return
		}

		_ = andonissueview.AddPage(&andonissueview.AddPageProps{
			Ctx:              ctx,
			Values:           r.Form,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
			AndonIssues:      andonIssues,
			Teams:            teams,
		}).Render(w)
		return
	}

	if err := h.andonIssueService.Create(
		r.Context(),
		model.NewAndonIssue{
			IssueName:          fd.IssueName,
			ParentID:           fd.ParentID,
			AssignedToTeam:     fd.AssignedToTeam,
			ResolvableByRaiser: fd.ResolvableByRaiser,
			WillStopProcess:    fd.WillStopProcess,
		},
		ctx.User.UserID,
	); err != nil {
		log.Println(err)
		http.Error(w, "Error creating andon issue", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/andon-issues", http.StatusSeeOther)
}

type addAndonIssueFormData struct {
	IssueName          string
	ParentID           *int
	AssignedToTeam     int
	ResolvableByRaiser bool
	WillStopProcess    bool
}

func (fd *addAndonIssueFormData) normalise() {
	fd.IssueName = strings.TrimSpace(fd.IssueName)
}

func (fd *addAndonIssueFormData) validate() validate.ValidationErrors {
	var ve validate.ValidationErrors = make(map[string][]string)

	validate.MinLength(&ve, "IssueName", fd.IssueName, 3)
	validate.MaxLength(&ve, "IssueName", fd.IssueName, 50)

	return ve
}

func (h *AndonIssueHandler) EditPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	andonIssueID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid andon issue ID", http.StatusBadRequest)
		return
	}

	andonIssue, err := h.andonIssueService.GetByID(r.Context(), andonIssueID)
	if err != nil {
		http.Error(w, "Error getting andon issue", http.StatusInternalServerError)
		return
	}
	if andonIssue == nil {
		http.Error(w, "Andon issue does not exist", http.StatusBadRequest)
		return
	}

	andonIssues, _, err := h.andonIssueService.List(
		r.Context(),
		model.ListAndonIssuesQuery{
			Page: 1, PageSize: 10000,
		})
	if err != nil {
		http.Error(w, "Error fetching andon issues", http.StatusInternalServerError)
		return
	}
	teams, _, err := h.teamService.List(r.Context(), model.ListTeamsQuery{
		Page: 1, PageSize: 10000,
	})
	if err != nil {
		http.Error(w, "Error fetching teams", http.StatusInternalServerError)
		return
	}

	_ = andonissueview.EditPage(&andonissueview.EditPageProps{
		Ctx:         ctx,
		AndonIssue:  *andonIssue,
		AndonIssues: andonIssues,
		Teams:       teams,
	}).Render(w)
}

func (h *AndonIssueHandler) Edit(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	andonIssueID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid andon issue ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var fd editAndonIssueFormData
	if err := appurl.Unmarshal(r.Form, &fd); err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd.normalise()

	validationErrors := fd.validate()

	if validationErrors == nil {
		validationErrors, err = h.andonIssueService.Update(
			r.Context(),
			andonIssueID,
			model.AndonIssueUpdate{
				IssueName:          fd.IssueName,
				ParentID:           fd.ParentID,
				IsArchived:         fd.IsArchived,
				AssignedToTeam:     fd.AssignedToTeam,
				ResolvableByRaiser: fd.ResolvableByRaiser,
				WillStopProcess:    fd.WillStopProcess,
			},
			ctx.User.UserID,
		)

		if err != nil {
			log.Println(err)
			http.Error(w, "Error updating andon issue", http.StatusInternalServerError)
			return
		}
	}

	if validationErrors != nil {
		andonIssues, _, err := h.andonIssueService.List(
			r.Context(),
			model.ListAndonIssuesQuery{
				Page: 1, PageSize: 10000,
			})
		if err != nil {
			http.Error(w, "Error fetching andon issues", http.StatusInternalServerError)
			return
		}
		teams, _, err := h.teamService.List(r.Context(), model.ListTeamsQuery{
			Page: 1, PageSize: 10000,
		})
		if err != nil {
			http.Error(w, "Error fetching teams", http.StatusInternalServerError)
			return
		}

		_ = andonissueview.EditPage(&andonissueview.EditPageProps{
			Ctx:              ctx,
			Values:           r.Form,
			ValidationErrors: *validationErrors,
			IsSubmission:     true,
			AndonIssues:      andonIssues,
			Teams:            teams,
		}).Render(w)
		return
	}

	http.Redirect(w, r, "/andon-issues", http.StatusSeeOther)
}

type editAndonIssueFormData struct {
	IssueName          string
	IsArchived         bool
	ParentID           *int
	AssignedToTeam     int
	ResolvableByRaiser bool
	WillStopProcess    bool
}

func (fd *editAndonIssueFormData) normalise() {
	fd.IssueName = strings.TrimSpace(fd.IssueName)
}

func (fd *editAndonIssueFormData) validate() *validate.ValidationErrors {
	var ve validate.ValidationErrors = make(map[string][]string)

	validate.MinLength(&ve, "IssueName", fd.IssueName, 3)
	validate.MaxLength(&ve, "IssueName", fd.IssueName, 50)

	if len(ve) == 0 {
		return nil
	}

	return &ve
}
