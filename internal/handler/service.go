package handler

import (
	"app/internal/model"
	"app/internal/service"
	"app/internal/views/serviceview"
	"app/pkg/apphmac"
	"app/pkg/appsort"
	"app/pkg/appurl"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type ServiceHandler struct {
	galleryService  service.GalleryService
	commentService  service.CommentService
	resourceService service.ResourceService
	servicesService service.ServicesService
	teamService     service.TeamService
	appHMAC         apphmac.AppHMAC
}

func NewServiceHandler(
	resourceService service.ResourceService,
	servicesService service.ServicesService,
	galleryService service.GalleryService,
	commentService service.CommentService,
	teamService service.TeamService,
	appHMAC apphmac.AppHMAC) *ServiceHandler {
	return &ServiceHandler{
		resourceService: resourceService,
		servicesService: servicesService,
		galleryService:  galleryService,
		commentService:  commentService,
		teamService:     teamService,
		appHMAC:         appHMAC,
	}
}

func (h *ServiceHandler) ServiceMetricsPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	var uv serviceMetricsPageURLVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		log.Println("error decoding url values:", err)
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	uv.normalise()

	sort := appsort.Sort{}
	err = sort.ParseQueryParam(model.ServiceMetric{}, uv.Sort)
	if err != nil {
		log.Println("error parsing resource sort:", err)
		http.Error(w, "Error parsing sort", http.StatusBadRequest)
		return
	}

	metrics, count, err := h.servicesService.GetServiceMetrics(r.Context(), uv.ShowArchived)
	if err != nil {
		log.Println("error listing service metrics:", err)
		http.Error(w, "Error listing service metrics", http.StatusInternalServerError)
		return
	}

	_ = serviceview.ServiceMetricsPage(&serviceview.ServiceMetricsPageProps{
		Ctx:          ctx,
		ShowArchived: uv.ShowArchived,
		Metrics:      metrics,
		MetricsCount: count,
		Sort:         sort,
		Page:         uv.Page,
		PageSize:     uv.PageSize,
	}).Render(w)
}

type serviceMetricsPageURLVals struct {
	Sort     string
	Page     int
	PageSize int

	ShowArchived bool
}

func (uv *serviceMetricsPageURLVals) normalise() {
	if uv.Page == 0 {
		uv.Page = 1
	}
	if uv.PageSize == 0 {
		uv.PageSize = 50
	}
}

func (h *ServiceHandler) ResourceServicingPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	var uv resourceServicingPageURLVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		log.Println("error decoding url values:", err)
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	uv.normalise()

	sort := appsort.Sort{}
	err = sort.ParseQueryParam(model.Resource{}, uv.Sort)
	if err != nil {
		log.Println("error parsing resource sort:", err)
		http.Error(w, "Error parsing sort", http.StatusBadRequest)
		return
	}

	resources, count, err := h.servicesService.GetResourceServiceMetricStatuses(r.Context(), model.ResourceServiceMetricStatusesQuery{
		ServiceOwnershipTeamIDs: uv.ServiceOwnershipTeamIDs,
		Page:                    uv.Page,
		PageSize:                uv.PageSize,
	})
	if err != nil {
		log.Println("error listing resources:", err)
		http.Error(w, "Error listing resources", http.StatusInternalServerError)
		return
	}

	teams, _, err := h.teamService.List(r.Context(), model.ListTeamsQuery{
		Page: 1, PageSize: 1000,
	})
	if err != nil {
		log.Println("error listing teams:", err)
		http.Error(w, "Error listing teams", http.StatusInternalServerError)
		return
	}

	selectedTeamIDs := make([]string, len(uv.ServiceOwnershipTeamIDs))
	for i, id := range uv.ServiceOwnershipTeamIDs {
		selectedTeamIDs[i] = strconv.Itoa(id)
	}

	_ = serviceview.ResourceServicingPage(&serviceview.ResourceServicingPageProps{
		Ctx:             ctx,
		Resources:       resources,
		Count:           count,
		Sort:            sort,
		Page:            uv.Page,
		PageSize:        uv.PageSize,
		Teams:           teams,
		SelectedTeamIDs: selectedTeamIDs,
	}).Render(w)
}

type resourceServicingPageURLVals struct {
	Sort     string
	Page     int
	PageSize int

	IsArchived              bool
	ServiceOwnershipTeamIDs []int
}

func (uv *resourceServicingPageURLVals) normalise() {
	if uv.Page == 0 {
		uv.Page = 1
	}
	if uv.PageSize == 0 {
		uv.PageSize = 50
	}
}

func (h *ServiceHandler) ServicesPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	var uv servicesPageURLVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		log.Println("error decoding url values:", err)
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	uv.normalise()

	sort := appsort.Sort{}
	err = sort.ParseQueryParam(model.ResourceService{}, uv.Sort)
	if err != nil {
		log.Println("error parsing resource sort:", err)
		http.Error(w, "Error parsing sort", http.StatusBadRequest)
		return
	}

	services, count, err := h.servicesService.ListServices(r.Context(), model.GetServicesQuery{
		Sort:     sort,
		Page:     uv.Page,
		PageSize: uv.PageSize,
	})
	if err != nil {
		log.Println("error listing resources:", err)
		http.Error(w, "Error listing resources", http.StatusInternalServerError)
		return
	}

	_ = serviceview.ServicesPage(&serviceview.ServicesPageProps{
		Ctx:      ctx,
		Services: services,
		Count:    count,
		Sort:     sort,
		Page:     uv.Page,
		PageSize: uv.PageSize,
	}).Render(w)
}

type servicesPageURLVals struct {
	Sort     string
	Page     int
	PageSize int

	IsArchived bool
}

func (uv *servicesPageURLVals) normalise() {
	if uv.Page == 0 {
		uv.Page = 1
	}
	if uv.PageSize == 0 {
		uv.PageSize = 50
	}
}

func (h *ServiceHandler) AddResourceServiceMetricPage(
	w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	type urlVals struct {
		Name         string
		Description  string
		IsCumulative bool
	}

	var uv urlVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding query params", http.StatusInternalServerError)
		return
	}

	_ = serviceview.AddServiceMetricPage(&serviceview.AddServiceMetricPageProps{
		Ctx:    ctx,
		Values: r.URL.Query(),
	}).Render(w)
}

func (h *ServiceHandler) AddResourceServiceMetric(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	type urlVals struct {
		Name         string
		Description  string
		IsCumulative bool
	}

	var uv urlVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		log.Println("error decoding query params:", err)
		http.Error(w, "Error decoding query params", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println("error parsing form:", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var fd addResourceServiceMetricFormData
	if err := appurl.Unmarshal(r.Form, &fd); err != nil {
		log.Println("error decoding form:", err)
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd.normalise()

	validationErrors := fd.validate()

	if len(validationErrors) > 0 {

		_ = serviceview.AddServiceMetricPage(
			&serviceview.AddServiceMetricPageProps{
				Ctx:              ctx,
				Values:           r.Form,
				ValidationErrors: validationErrors,
				IsSubmission:     true,
			}).Render(w)
		return
	}

	err = h.servicesService.CreateResourceServiceMetric(
		r.Context(),
		model.NewServiceMetric{
			Name:         fd.Name,
			Description:  fd.Description,
			IsCumulative: fd.IsCumulative,
		})
	if err != nil {
		log.Println("error creating resource metric:", err)
		http.Error(w, "Error creating resource metric", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/services/metrics", http.StatusSeeOther)
}

func (h *ServiceHandler) EditResourceServiceMetricPage(
	w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	metricID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println("invalid metric id:", err)
		http.Error(w, "Invalid metric id", http.StatusBadRequest)
		return
	}

	metric, err := h.servicesService.GetResourceServiceMetricByID(r.Context(), metricID)
	if err != nil {
		log.Println("error fetching resource metric:", err)
		http.Error(w, "Error fetching resource metric", http.StatusInternalServerError)
		return
	}
	if metric == nil {
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}

	_ = serviceview.EditServiceMetricPage(&serviceview.EditServiceMetricPageProps{
		Ctx:    ctx,
		Metric: *metric,
		Values: r.URL.Query(),
	}).Render(w)
}

func (h *ServiceHandler) EditResourceServiceMetric(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	metricID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println("invalid metric id:", err)
		http.Error(w, "Invalid metric id", http.StatusBadRequest)
		return
	}

	metric, err := h.servicesService.GetResourceServiceMetricByID(r.Context(), metricID)
	if err != nil {
		log.Println("error fetching resource metric:", err)
		http.Error(w, "Error fetching resource metric", http.StatusInternalServerError)
		return
	}
	if metric == nil {
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println("error parsing form:", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var fd editResourceServiceMetricFormData
	if err := appurl.Unmarshal(r.Form, &fd); err != nil {
		log.Println("error decoding form:", err)
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd.normalise()

	validationErrors := fd.validate()

	if len(validationErrors) > 0 {
		_ = serviceview.EditServiceMetricPage(
			&serviceview.EditServiceMetricPageProps{
				Ctx:              ctx,
				Metric:           *metric,
				Values:           r.Form,
				ValidationErrors: validationErrors,
				IsSubmission:     true,
			}).Render(w)
		return
	}

	err = h.servicesService.UpdateResourceServiceMetric(
		r.Context(),
		model.UpdateResourceServiceMetric{
			ServiceMetricID: metricID,
			Name:            fd.Name,
			Description:     fd.Description,
			IsCumulative:    fd.IsCumulative,
			IsArchived:      fd.IsArchived,
		},
	)
	if err != nil {
		log.Println("error updating resource metric:", err)
		http.Error(w, "Error updating resource metric", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/services/metrics", http.StatusSeeOther)
}

type addResourceServiceMetricFormData struct {
	Name         string
	Description  string
	IsCumulative bool
}

func (fd *addResourceServiceMetricFormData) normalise() {
	fd.Name = strings.TrimSpace(fd.Name)
	fd.Description = strings.TrimSpace(fd.Description)
}

func (fd *addResourceServiceMetricFormData) validate() validate.ValidationErrors {
	var ve validate.ValidationErrors = make(map[string][]string)

	if fd.Name == "" {
		ve.Add("Name", "is required")
	}

	if fd.Description == "" {
		ve.Add("Description", "is required")
	}

	return ve
}

type editResourceServiceMetricFormData struct {
	Name         string
	Description  string
	IsCumulative bool
	IsArchived   bool
}

func (fd *editResourceServiceMetricFormData) normalise() {
	fd.Name = strings.TrimSpace(fd.Name)
	fd.Description = strings.TrimSpace(fd.Description)
}

func (fd *editResourceServiceMetricFormData) validate() validate.ValidationErrors {
	var ve validate.ValidationErrors = make(map[string][]string)

	if fd.Name == "" {
		ve.Add("Name", "is required")
	}

	if fd.Description == "" {
		ve.Add("Description", "is required")
	}

	return ve
}

func (h *ServiceHandler) DeleteResourceServiceMetric(w http.ResponseWriter, r *http.Request) {

	metricID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println("invalid metric id:", err)
		http.Error(w, "Invalid metric id", http.StatusBadRequest)
		return
	}

	err = h.servicesService.DeleteResourceServiceMetric(r.Context(), metricID)
	if err != nil {
		log.Println("error deleting resource metric:", err)
		http.Error(w, "Error deleting resource metric", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/services/metrics", http.StatusSeeOther)
}

type editResourceServiceFormData struct {
	ResourceID int
	Notes      string
}

func (fd *editResourceServiceFormData) normalise() {
	fd.Notes = strings.TrimSpace(fd.Notes)
}

func (fd *editResourceServiceFormData) validate() validate.ValidationErrors {
	var ve validate.ValidationErrors = make(map[string][]string)

	return ve
}

func (h *ServiceHandler) UpdateResourceServiceNotes(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	serviceID, err := strconv.Atoi(r.PathValue("serviceID"))
	if err != nil {
		log.Println("invalid resource service id:", err)
		http.Error(w, "Invalid resource service id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println("error parsing form:", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var fd editResourceServiceFormData
	if err := appurl.Unmarshal(r.Form, &fd); err != nil {
		log.Println("error decoding form:", err)
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd.normalise()

	validationErrors := fd.validate()

	if len(validationErrors) > 0 {
		_ = serviceview.ResourceServicePage(
			&serviceview.ResourceServicePageProps{
				Ctx:              ctx,
				Values:           r.Form,
				ValidationErrors: validationErrors,
				IsSubmission:     true,
			}).Render(w)
		return
	}

	err = h.servicesService.UpdateResourceService(
		r.Context(),
		model.UpdateResourceService{
			ResourceServiceID: serviceID,
			Notes:             fd.Notes,
		}, ctx.User.UserID)
	if err != nil {
		if errors.Is(err, service.ErrResourceServiceNotFound) {
			http.Error(w, "Resource service not found", http.StatusNotFound)
			return
		}
		log.Println("error updating resource service:", err)
		http.Error(w, "Error updating resource service", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/services/%d", serviceID), http.StatusSeeOther)
}

func (h *ServiceHandler) ResourceServicePage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	userID := ctx.User.UserID

	canUserEdit := ctx.User.Permissions.UserAdmin.Access

	serviceID, err := strconv.Atoi(r.PathValue("serviceID"))
	if err != nil {
		http.Error(w, "Invalid resource service id", http.StatusBadRequest)
		return
	}

	resourceService, err := h.servicesService.GetResourceServiceByID(r.Context(), serviceID)
	if err != nil {
		if errors.Is(err, service.ErrResourceServiceNotFound) {
			http.Error(w, "Resource service not found", http.StatusNotFound)
			return
		}
		log.Println("error fetching resource service:", err)
		http.Error(w, "Error fetching resource service", http.StatusInternalServerError)
		return
	}

	lastService, err := h.servicesService.GetLastServiceForResource(
		r.Context(),
		resourceService.ResourceID,
		resourceService.ResourceServiceID,
		resourceService.StartedAt,
	)
	if err != nil {
		log.Println("error fetching last service for resource:", err)
		http.Error(w, "Error fetching last service for resource", http.StatusInternalServerError)
		return
	}

	changelog, err := h.servicesService.GetServiceChangelog(
		r.Context(),
		serviceID,
	)
	if err != nil {
		log.Println("error fetching resource service changelog:", err)
		http.Error(w, "Error fetching resource service changelog", http.StatusInternalServerError)
		return
	}

	if canUserEdit {
		resourceService.GalleryURL = h.galleryService.GenerateEditTempURL(resourceService.GalleryID, true)
	} else {
		resourceService.GalleryURL = h.galleryService.GenerateTempURL(resourceService.GalleryID, canUserEdit)
	}

	galleryImgURLs, err := h.galleryService.GetGalleryImgURLs(r.Context(), resourceService.GalleryID)
	if err != nil {
		log.Println("error fetching service gallery:", err)
		http.Error(w, "Error fetching service gallery", http.StatusInternalServerError)
		return
	}

	serviceComments, err := h.commentService.GetComments(r.Context(), resourceService.CommentThreadID, userID)
	if err != nil {
		log.Println("error fetching service comments:", err)
		http.Error(w, "Error fetching service comments", http.StatusInternalServerError)
		return
	}

	// Build a JSON envelope for adding a comment to this andon's thread, valid for 24 hours
	permissions := []string{"view"}
	if canUserEdit {
		permissions = append(permissions, "add")
	}

	commentPayload := apphmac.Payload{
		Entity:      "comment_thread",
		EntityID:    fmt.Sprintf("%d", resourceService.CommentThreadID),
		Permissions: permissions,
		Expires:     time.Now().Add(24 * time.Hour).Unix(),
	}
	commentEnvelope := h.appHMAC.CreateEnvelope(
		commentPayload,
	)

	_ = serviceview.ResourceServicePage(&serviceview.ResourceServicePageProps{
		Ctx:                     ctx,
		ResourceService:         *resourceService,
		LastResourceService:     lastService,
		GalleryImageURLs:        galleryImgURLs,
		ResourceServiceComments: serviceComments,
		ServiceChangelog:        changelog,
		CommentHMACEnvelope:     commentEnvelope,
	}).Render(w)
}

func (h *ServiceHandler) UpdateResourceService(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	resourceID, err := strconv.Atoi(r.PathValue("resourceID"))
	if err != nil {
		http.Error(w, "Invalid resource id", http.StatusBadRequest)
		return
	}
	action := r.PathValue("action")

	serviceID, err := strconv.Atoi(r.PathValue("serviceID"))
	if err != nil {
		http.Error(w, "Invalid resource service id", http.StatusBadRequest)
		return
	}

	resource, err := h.resourceService.GetResourceByID(r.Context(), resourceID)
	if err != nil {
		log.Println("error fetching resource:", err)
		http.Error(w, "Error fetching resource", http.StatusInternalServerError)
		return
	}

	if resource == nil {
		http.Error(w, "Resource not found", http.StatusNotFound)
		return
	}

	err = h.servicesService.UpdateService(
		r.Context(),
		resourceID,
		serviceID,
		action,
		ctx.User.UserID,
	)
	if err != nil {
		if errors.Is(err, service.ErrResourceServiceNotFound) {
			http.Error(w, "Resource service not found", http.StatusNotFound)
			return
		}
		log.Println(err)
		http.Error(w, "Error updating resource service", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ServiceHandler) AddServiceSchedulePage(
	w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	resourceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println("invalid resource id:", err)
		http.Error(w, "Invalid resource id", http.StatusBadRequest)
		return
	}

	resource, err := h.resourceService.GetResourceByID(r.Context(), resourceID)
	if err != nil {
		log.Println("error fetching resource:", err)
		http.Error(w, "Error fetching resource", http.StatusInternalServerError)
		return
	}
	if resource == nil {
		http.Error(w, "Resource not found", http.StatusNotFound)
		return
	}

	metrics, _, err := h.servicesService.GetServiceMetrics(r.Context(), false)
	if err != nil {
		log.Println("error fetching service metrics:", err)
		http.Error(w, "Error fetching service metrics", http.StatusInternalServerError)
		return
	}

	_ = serviceview.AddServiceSchedulePage(&serviceview.AddResourceServiceSchedulePageProps{
		Ctx:            ctx,
		Values:         r.URL.Query(),
		Resource:       *resource,
		ServiceMetrics: metrics,
	}).Render(w)
}

func (h *ServiceHandler) AddServiceSchedule(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	resourceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println("invalid resource id:", err)
		http.Error(w, "Invalid resource id", http.StatusBadRequest)
		return
	}

	resource, err := h.resourceService.GetResourceByID(r.Context(), resourceID)
	if err != nil {
		log.Println("error fetching resource:", err)
		http.Error(w, "Error fetching resource", http.StatusInternalServerError)
		return
	}
	if resource == nil || resource.IsArchived {
		http.Error(w, "Resource not available", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println("error parsing form:", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var fd addServiceScheduleFormData
	if err := appurl.Unmarshal(r.Form, &fd); err != nil {
		log.Println("error decoding form:", err)
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd.normalise()

	validationErrors := fd.validate()

	if len(validationErrors) > 0 {
		metrics, _, err := h.servicesService.GetServiceMetrics(r.Context(), false)
		if err != nil {
			log.Println("error fetching service metrics:", err)
			http.Error(w, "Error fetching service metrics", http.StatusInternalServerError)
			return
		}

		_ = serviceview.AddServiceSchedulePage(
			&serviceview.AddResourceServiceSchedulePageProps{
				Ctx:              ctx,
				Values:           r.Form,
				ValidationErrors: validationErrors,
				IsSubmission:     true,
				Resource:         *resource,
				ServiceMetrics:   metrics,
			}).Render(w)
		return
	}

	err = h.servicesService.CreateResourceServiceSchedule(
		r.Context(),
		model.NewResourceServiceSchedule{
			ResourceID:              resourceID,
			ResourceServiceMetricID: fd.ServiceMetricID,
			Threshold:               fd.Threshold,
		})
	if err != nil {
		log.Println("error creating resource service schedule:", err)
		http.Error(w, "Error creating resource service schedule", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/resources/%d", resourceID), http.StatusSeeOther)
}

func (h *ServiceHandler) ArchiveServiceSchedule(w http.ResponseWriter, r *http.Request) {

	resourceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println("invalid resource id:", err)
		http.Error(w, "Invalid resource id", http.StatusBadRequest)
		return
	}

	scheduleID, err := strconv.Atoi(r.PathValue("scheduleID"))
	if err != nil {
		log.Println("invalid schedule id:", err)
		http.Error(w, "Invalid schedule id", http.StatusBadRequest)
		return
	}

	err = h.servicesService.ArchiveResourceServiceSchedule(r.Context(), resourceID, scheduleID)
	if err != nil {
		log.Println("error archiving service schedule:", err)
		http.Error(w, "Error archiving service schedule", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/resources/%d", resourceID), http.StatusSeeOther)
}

type addServiceScheduleFormData struct {
	ServiceMetricID int
	Threshold       decimal.Decimal
}

func (fd *addServiceScheduleFormData) normalise() {
}

func (fd *addServiceScheduleFormData) validate() validate.ValidationErrors {
	var ve validate.ValidationErrors = make(map[string][]string)

	if fd.ServiceMetricID == 0 {
		ve.Add("ServiceMetricID", "must be selected")
	}
	validate.DecimalGT(&ve, "Threshold", fd.Threshold, decimal.Zero)

	return ve
}
