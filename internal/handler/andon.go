package handler

import (
	"app/internal/model"
	"app/internal/service"
	"app/internal/views/andonview"
	"app/pkg/appsort"
	"app/pkg/appurl"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AndonHandler struct {
	andonService      service.AndonService
	andonIssueService service.AndonIssueService
	commentService    service.CommentService
	teamService       service.TeamService
	fileService       service.FileService
}

func NewAndonHandler(
	andonService service.AndonService,
	andonIssueService service.AndonIssueService,
	commentService service.CommentService,
	teamService service.TeamService,
	fileService service.FileService,
) *AndonHandler {
	return &AndonHandler{
		andonService:      andonService,
		andonIssueService: andonIssueService,
		commentService:    commentService,
		teamService:       teamService,
		fileService:       fileService,
	}
}

func (h *AndonHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var uv andonEventsHomePageUrlVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	uv.normalise()

	sort := appsort.Sort{}
	err = sort.ParseQueryParam(model.AndonEvent{}, uv.Sort)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Error parsing sort: %v", err), http.StatusBadRequest)
		return
	}

	outstandingAndons, outstandingAndonsCount, _, err := h.andonService.ListAndons(r.Context(),
		model.ListAndonQuery{
			Sort:     sort,
			Page:     uv.Page,
			PageSize: uv.PageSize,

			Teams:            uv.AndonTeams,
			Statuses:         []string{"Outstanding"},
			OrderBy:          "raised_at",
			OrderByDirection: "asc",
		}, ctx.User.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error listing outstanding andons", http.StatusInternalServerError)
		return
	}

	acknowledgedAndons, acknowledgedAndonsCount, _, err := h.andonService.ListAndons(
		r.Context(),
		model.ListAndonQuery{
			Sort:     sort,
			Page:     uv.Page,
			PageSize: uv.PageSize,

			Teams:            uv.AndonTeams,
			Statuses:         []string{"Acknowledged"},
			OrderBy:          "acknowledged_at",
			OrderByDirection: "asc",
		},
		ctx.User.UserID,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error listing acknowledged andons", http.StatusInternalServerError)
		return
	}

	teams, _, err := h.teamService.List(r.Context(), model.ListTeamsQuery{
		Page: 1, PageSize: 10000,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "Error fetching teams", http.StatusInternalServerError)
		return
	}

	_ = andonview.HomePage(&andonview.HomePageProps{
		Ctx:                     ctx,
		OutstandingAndons:       outstandingAndons,
		AcknowledgedAndons:      acknowledgedAndons,
		NewAndonsCount:          outstandingAndonsCount,
		AcknowledgedAndonsCount: acknowledgedAndonsCount,
		Teams:                   teams,
		SelectedTeams:           uv.AndonTeams,
		Sort:                    sort,
		Page:                    uv.Page,
		PageSize:                uv.PageSize,
	}).Render(w)
}

type andonEventsHomePageUrlVals struct {
	Sort     string
	Page     int
	PageSize int

	Status     string
	AndonTeams []string
}

func (uv *andonEventsHomePageUrlVals) normalise() {
	if uv.Page == 0 {
		uv.Page = 1
	}
	if uv.PageSize == 0 {
		uv.PageSize = 50
	}
}

type andonAllEventsUrlVals struct {
	Sort     string
	Page     int
	PageSize int

	StartDate                *time.Time
	EndDate                  *time.Time
	IssueIn                  []string
	SeverityIn               []string
	TeamIn                   []string
	LocationIn               []string
	StatusIn                 []string
	RaisedByUsernameIn       []string
	AcknowledgedByUsernameIn []string
	ResolvedByUsernameIn     []string
}

func (uv *andonAllEventsUrlVals) normalise() {
	if uv.Page == 0 {
		uv.Page = 1
	}
	if uv.PageSize == 0 {
		uv.PageSize = 50
	}
}

func (h *AndonHandler) AllAndonsPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var uv andonAllEventsUrlVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	uv.normalise()

	sort := appsort.Sort{}
	err = sort.ParseQueryParam(model.AndonEvent{}, uv.Sort)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Error parsing sort: %v", err), http.StatusBadRequest)
		return
	}

	andons, count, filters, err := h.andonService.ListAndons(
		r.Context(),
		model.ListAndonQuery{
			Sort:     sort,
			Page:     uv.Page,
			PageSize: uv.PageSize,

			StartDate:              uv.StartDate,
			EndDate:                uv.EndDate,
			Issues:                 uv.IssueIn,
			Serverities:            uv.SeverityIn,
			Teams:                  uv.TeamIn,
			Locations:              uv.LocationIn,
			Statuses:               uv.StatusIn,
			RaisedByUsername:       uv.RaisedByUsernameIn,
			AcknowledgedByUsername: uv.AcknowledgedByUsernameIn,
			ResolvedByUsername:     uv.ResolvedByUsernameIn,
		},
		ctx.User.UserID,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error listing andons", http.StatusInternalServerError)
		return
	}

	_ = andonview.AllAndonsPage(&andonview.AllAndonsPageProps{
		Ctx:              ctx,
		Andons:           andons,
		AndonsCount:      count,
		AvailableFilters: filters,
		Sort:             sort,
		Page:             uv.Page,
		PageSize:         uv.PageSize,
		Filters: model.AndonFilters{
			StartDate:              uv.StartDate,
			EndDate:                uv.EndDate,
			Issues:                 uv.IssueIn,
			Severities:             uv.SeverityIn,
			Teams:                  uv.TeamIn,
			Locations:              uv.LocationIn,
			Statuses:               uv.StatusIn,
			RaisedByUsername:       uv.RaisedByUsernameIn,
			AcknowledgedByUsername: uv.AcknowledgedByUsernameIn,
			ResolvedByUsername:     uv.ResolvedByUsernameIn,
		},
	}).Render(w)
}

func (h *AndonHandler) AddPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	type urlVals struct {
		IssueOrGroupID   int
		Source           string
		LinkedEntityID   int
		LinkedEntityType string
		Location         string
		ReturnTo         string
	}

	var uv urlVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding query params", http.StatusInternalServerError)
		return
	}

	if uv.IssueOrGroupID != 0 {
		issueNodes, err := h.andonIssueService.GetIssueHierarchy(r.Context(), uv.IssueOrGroupID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to get issue hierarchy", http.StatusInternalServerError)
			return
		}

		cleanValues := r.URL.Query()

		cleanValues.Del("IssueOrGroupID")

		for i, node := range issueNodes {
			paramKey := fmt.Sprintf("Node[%d]", i)
			cleanValues.Set(paramKey, strconv.Itoa(node))
		}

		cleanURL := *r.URL
		cleanURL.RawQuery = cleanValues.Encode()

		http.Redirect(w, r, cleanURL.String(), http.StatusFound)
		return
	}

	nodes := []int{}
	for i := 0; ; i++ {
		nodeStr := r.URL.Query().Get(fmt.Sprintf("Node[%d]", i))
		if nodeStr == "" {
			break
		}
		nodeID, err := strconv.Atoi(nodeStr)
		if err != nil {
			break
		}
		nodes = append(nodes, nodeID)
	}

	andonIssues, _, err := h.andonIssueService.ListIssuesAndGroups(
		r.Context(),
		model.ListAndonIssuesQuery{
			Page: 1, PageSize: 10000,
		})
	if err != nil {
		log.Println(err)
		http.Error(w, "Error fetching andon issues", http.StatusInternalServerError)
		return
	}

	teams, _, err := h.teamService.List(r.Context(), model.ListTeamsQuery{
		Page: 1, PageSize: 10000,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "Error fetching teams", http.StatusInternalServerError)
		return
	}

	_ = andonview.AddPage(&andonview.AddPageProps{
		Ctx:          ctx,
		Values:       r.URL.Query(),
		AndonIssues:  andonIssues,
		Teams:        teams,
		SelectedPath: nodes,
	}).Render(w)
}

func (h *AndonHandler) Add(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	type urlVals struct {
		Source   string
		ReturnTo string
	}

	var uv urlVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding query params", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var fd addAndonEventFormData
	if err := appurl.Unmarshal(r.Form, &fd); err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd.normalise()

	validationErrors := fd.validate()

	nodes := []int{}
	for i := 0; ; i++ {
		nodeStr := r.URL.Query().Get(fmt.Sprintf("Node[%d]", i))
		if nodeStr == "" {
			break
		}
		nodeID, err := strconv.Atoi(nodeStr)
		if err != nil {
			break
		}
		nodes = append(nodes, nodeID)
	}

	if len(validationErrors) > 0 {
		andonIssues, _, err := h.andonIssueService.ListIssuesAndGroups(
			r.Context(),
			model.ListAndonIssuesQuery{
				Page: 1, PageSize: 10000,
			})
		if err != nil {
			log.Println(err)
			http.Error(w, "Error fetching andon issues", http.StatusInternalServerError)
			return
		}

		teams, _, err := h.teamService.List(r.Context(), model.ListTeamsQuery{
			Page: 1, PageSize: 10000,
		})
		if err != nil {
			log.Println(err)
			http.Error(w, "Error fetching teams", http.StatusInternalServerError)
			return
		}

		_ = andonview.AddPage(&andonview.AddPageProps{
			Ctx:              ctx,
			Values:           r.Form,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
			AndonIssues:      andonIssues,
			Teams:            teams,
			SelectedPath:     nodes,
		}).Render(w)
		return
	}

	if err := h.andonService.CreateAndonEvent(
		r.Context(),
		model.NewAndonEvent{
			IssueDescription: fd.IssueDescription,
			IssueID:          fd.IssueID,
			Location:         fd.Location,
			Source:           uv.Source,
			LinkedEntityID:   fd.LinkedEntityID,
			LinkedEntityType: fd.LinkedEntityType,
		},
		ctx.User.UserID,
	); err != nil {
		http.Error(w, "Error creating andon", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/andons", http.StatusSeeOther)
}

func (h *AndonHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	entity := "Andon"
	entityIDStr := r.PathValue("entityID")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var fd addAndonCommentFormData
	if err := json.NewDecoder(r.Body).Decode(&fd); err != nil {
		log.Println(err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	fd.normalise()

	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		http.Error(w, "Invalid entity Id", http.StatusBadRequest)
		return
	}

	commentID, err := h.commentService.CreateComment(
		r.Context(),
		&model.NewComment{
			Comment:  fd.Comment,
			Entity:   entity,
			EntityID: entityID,
		},
		ctx.User.UserID,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error creating andon", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"commentId": commentID,
	})
}

func (h *AndonHandler) AddAttachment(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	entityIDStr := r.PathValue("commentID")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		http.Error(w, "Invalid entity Id", http.StatusBadRequest)
		return
	}

	var fd addFileFormData
	if err := json.NewDecoder(r.Body).Decode(&fd); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	file, signedURL, err := h.fileService.CreateFile(
		r.Context(),
		&model.File{
			Filename:    fd.Filename,
			ContentType: fd.ContentType,
			SizeBytes:   fd.SizeBytes,
			Entity:      "Comment",
			EntityID:    entityID,
		},
		ctx.User.UserID,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error adding file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"fileId":    file.FileID,
		"signedUrl": signedURL,
	})

}

func (h *AndonHandler) AndonUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	andonEventID, _ := strconv.Atoi(r.PathValue("andonID"))
	andonAction := r.PathValue("action")

	err := h.andonService.UpdateAndonEvent(
		r.Context(),
		andonEventID,
		andonAction,
		ctx.User.UserID,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error updating andon", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/andons", http.StatusSeeOther)
}

func (h *AndonHandler) AndonDetailsPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	andonID, _ := strconv.Atoi(r.PathValue("andonID"))

	andonEvent, err := h.andonService.GetAndonEventByID(r.Context(), andonID, ctx.User.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get andon event", http.StatusInternalServerError)
		return
	}

	changes, err := h.andonService.GetAndonByID(
		r.Context(),
		andonID,
		ctx.User.UserID,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error fetching andon details", http.StatusInternalServerError)
		return
	}

	comments, err := h.commentService.GetComments(r.Context(), "Andon", andonID, ctx.User.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error fetching andon comments", http.StatusInternalServerError)
		return
	}

	_ = andonview.AndonDetailsPage(&andonview.AndonDetailsPageProps{
		Ctx:           ctx,
		Values:        r.Form,
		IsSubmission:  true,
		AndonID:       andonID,
		AndonEvent:    *andonEvent,
		AndonChanges:  changes,
		AndonComments: comments,
	}).Render(w)
}

type addAndonEventFormData struct {
	IssueDescription string
	IssueID          int
	AssignedTeam     string
	Location         string
	LinkedEntityID   int
	LinkedEntityType string
}

func (fd *addAndonEventFormData) normalise() {
	fd.IssueDescription = strings.TrimSpace(fd.IssueDescription)
}

func (fd *addAndonEventFormData) validate() validate.ValidationErrors {
	var ve validate.ValidationErrors = make(map[string][]string)

	if fd.IssueDescription == "" {
		ve.Add("IssueDescription", "is required")
	}

	if fd.IssueID == 0 {
		ve.Add("IssueID", "is required")
	}

	if fd.Location == "" {
		ve.Add("Location", "is required")
	}

	if fd.AssignedTeam == "" {
		ve.Add("AssignedTeam", "for the issue is not present")
	}

	return ve
}
