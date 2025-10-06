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
	fileService       service.FileService
	galleryService    service.GalleryService
	teamService       service.TeamService
}

func NewAndonHandler(
	andonService service.AndonService,
	andonIssueService service.AndonIssueService,
	commentService service.CommentService,
	fileService service.FileService,
	galleryService service.GalleryService,
	teamService service.TeamService,
) *AndonHandler {
	return &AndonHandler{
		andonService:      andonService,
		andonIssueService: andonIssueService,
		commentService:    commentService,
		fileService:       fileService,
		galleryService:    galleryService,
		teamService:       teamService,
	}
}

func (h *AndonHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var uv andonsHomePageUrlVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	newSort := appsort.Sort{}
	err = newSort.ParseQueryParam(model.Andon{}, uv.NewSort)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Error parsing outstanding sort: %v", err), http.StatusBadRequest)
		return
	}
	wipSort := appsort.Sort{}
	err = wipSort.ParseQueryParam(model.Andon{}, uv.WIPSort)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Error parsing WIP sort: %v", err), http.StatusBadRequest)
		return
	}

	falseBool := false
	trueBool := true

	newAndons, newAndonsCount, _, err := h.andonService.ListAndons(r.Context(),
		model.ListAndonQuery{
			Sort:                 newSort,
			DefaultSortField:     "raised_at",
			DefaultSortDirection: appsort.DirectionAsc,
			Page:                 1,
			PageSize:             10000, // reasonable limit

			IsAcknowledged: &falseBool,
			IsOpen:         &trueBool,
			TeamIn:         uv.AndonTeams,
		}, ctx.User.UserID)
	if err != nil {
		log.Println("error listing outstanding andons:", err)
		http.Error(w, "Error listing outstanding andons", http.StatusInternalServerError)
		return
	}

	wipAndons, wipAndonsCount, _, err := h.andonService.ListAndons(
		r.Context(),
		model.ListAndonQuery{
			Sort:                 wipSort,
			DefaultSortField:     "raised_at",
			DefaultSortDirection: appsort.DirectionAsc,
			Page:                 1,
			PageSize:             10000, // reasonable limit

			IsAcknowledged: &trueBool,
			IsOpen:         &trueBool,
			TeamIn:         uv.AndonTeams,
		},
		ctx.User.UserID,
	)
	if err != nil {
		log.Println("error listing work in progress andons:", err)
		http.Error(w, "Error listing work in progress andons", http.StatusInternalServerError)
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
		Ctx:            ctx,
		NewAndons:      newAndons,
		WIPAndons:      wipAndons,
		NewAndonsCount: newAndonsCount,
		WIPAndonsCount: wipAndonsCount,
		Teams:          teams,
		SelectedTeams:  uv.AndonTeams,
		NewSort:        newSort,
		WIPSort:        wipSort,
		ReturnTo:       uv.ReturnTo,
	}).Render(w)
}

type andonsHomePageUrlVals struct {
	NewSort    string
	WIPSort    string
	AndonTeams []string
	ReturnTo   string
}

type allAndonsURLVals struct {
	Sort     string
	Page     int
	PageSize int

	StartDate                *time.Time
	EndDate                  *time.Time
	LocationIn               []string
	IssueIn                  []string
	SeverityIn               []string
	StatusIn                 []string
	TeamIn                   []string
	RaisedByUsernameIn       []string
	AcknowledgedByUsernameIn []string
	ResolvedByUsernameIn     []string
}

func (uv *allAndonsURLVals) normalise() {
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

	var uv allAndonsURLVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	uv.normalise()

	sort := appsort.Sort{}
	err = sort.ParseQueryParam(model.Andon{}, uv.Sort)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Error parsing sort: %v", err), http.StatusBadRequest)
		return
	}

	andons, count, availableFilters, err := h.andonService.ListAndons(
		r.Context(),
		model.ListAndonQuery{
			Sort:                 sort,
			DefaultSortField:     "andon_id",
			DefaultSortDirection: appsort.DirectionDesc,
			Page:                 uv.Page,
			PageSize:             uv.PageSize,

			StartDate:                uv.StartDate,
			EndDate:                  uv.EndDate,
			LocationIn:               uv.LocationIn,
			IssueIn:                  uv.IssueIn,
			SeverityIn:               uv.SeverityIn,
			StatusIn:                 uv.StatusIn,
			TeamIn:                   uv.TeamIn,
			RaisedByUsernameIn:       uv.RaisedByUsernameIn,
			AcknowledgedByUsernameIn: uv.AcknowledgedByUsernameIn,
			ResolvedByUsernameIn:     uv.ResolvedByUsernameIn,
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
		AvailableFilters: availableFilters,
		Sort:             sort,
		Page:             uv.Page,
		PageSize:         uv.PageSize,
		ActiveFilters: model.AndonFilters{
			StartDate:                uv.StartDate,
			EndDate:                  uv.EndDate,
			IssueIn:                  uv.IssueIn,
			SeverityIn:               uv.SeverityIn,
			StatusIn:                 uv.StatusIn,
			TeamIn:                   uv.TeamIn,
			LocationIn:               uv.LocationIn,
			RaisedByUsernameIn:       uv.RaisedByUsernameIn,
			AcknowledgedByUsernameIn: uv.AcknowledgedByUsernameIn,
			ResolvedByUsernameIn:     uv.ResolvedByUsernameIn,
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

	var fd addAndonFormData
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
		model.NewAndon{
			Description: fd.Description,
			IssueID:     fd.IssueID,
			Location:    fd.Location,
			Source:      fd.Source,
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
	json.NewEncoder(w).Encode(map[string]any{
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
	json.NewEncoder(w).Encode(map[string]any{
		"fileId":    file.FileID,
		"signedUrl": signedURL,
	})
}

func (h *AndonHandler) UpdateAndon(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	andonEventID, _ := strconv.Atoi(r.PathValue("andonID"))
	andonAction := r.PathValue("action")

	err := h.andonService.UpdateAndon(
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

type addAndonFormData struct {
	Description string
	IssueID     int
	Location    string
	Source      string
}

func (fd *addAndonFormData) normalise() {
	fd.Description = strings.TrimSpace(fd.Description)
}

func (fd *addAndonFormData) validate() validate.ValidationErrors {
	var ve validate.ValidationErrors = make(map[string][]string)

	if fd.IssueID == 0 {
		ve.Add("IssueID", "is required")
	}

	if fd.Location == "" {
		ve.Add("Location", "is required")
	}

	return ve
}

func (h *AndonHandler) AndonPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	andonID, _ := strconv.Atoi(r.PathValue("andonID"))

	type urlVals struct {
		ReturnTo string
	}

	var uv urlVals
	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding query params", http.StatusInternalServerError)
		return
	}

	andon, err := h.andonService.GetAndonByID(r.Context(), andonID, ctx.User.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error fetching andon", http.StatusInternalServerError)
		return
	}

	changelog, err := h.andonService.GetAndonChangelog(
		r.Context(),
		andonID,
		ctx.User.UserID,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error fetching andon changelog", http.StatusInternalServerError)
		return
	}

	galleryImgURLs, err := h.galleryService.GetGalleryImgURLs(r.Context(), andon.GalleryID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error fetching andon gallery", http.StatusInternalServerError)
		return
	}

	var galleryURL string
	if len(galleryImgURLs) == 0 && andon.CanUserEdit {
		galleryURL = h.galleryService.GenerateEditTempURL(andon.GalleryID, true)
	} else {
		galleryURL = h.galleryService.GenerateTempURL(andon.GalleryID, andon.CanUserEdit)
	}

	comments, err := h.commentService.GetComments(r.Context(), "Andon", andonID, ctx.User.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error fetching andon comments", http.StatusInternalServerError)
		return
	}

	_ = andonview.AndonPage(&andonview.AndonPageProps{
		Ctx:              ctx,
		Values:           r.Form,
		IsSubmission:     true,
		AndonID:          andonID,
		Andon:            *andon,
		AndonChangelog:   changelog,
		AndonComments:    comments,
		GalleryURL:       galleryURL,
		GalleryImageURLs: galleryImgURLs,
		ReturnTo:         uv.ReturnTo,
	}).Render(w)
}
