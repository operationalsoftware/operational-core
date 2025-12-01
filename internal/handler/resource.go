package handler

import (
	"app/internal/model"
	"app/internal/service"
	"app/internal/views/resourceview"
	"app/pkg/apphmac"
	"app/pkg/appsort"
	"app/pkg/appurl"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

type ResourceHandler struct {
	galleryService  service.GalleryService
	commentService  service.CommentService
	resourceService service.ResourceService
	servicesService service.ServicesService
	teamService     service.TeamService
	appHMAC         apphmac.AppHMAC
}

func NewResourceHandler(
	resourceService service.ResourceService,
	servicesService service.ServicesService,
	galleryService service.GalleryService,
	commentService service.CommentService,
	teamService service.TeamService,
	appHMAC apphmac.AppHMAC) *ResourceHandler {
	return &ResourceHandler{
		resourceService: resourceService,
		servicesService: servicesService,
		galleryService:  galleryService,
		commentService:  commentService,
		teamService:     teamService,
		appHMAC:         appHMAC,
	}
}

func (h *ResourceHandler) ResourcesPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	var uv resourceHomePageURLVals

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

	resources, count, err := h.resourceService.GetResources(r.Context(), model.GetResourcesQuery{
		Sort:     sort,
		Page:     uv.Page,
		PageSize: uv.PageSize,

		IsArchived: uv.IsArchived,
	})
	if err != nil {
		log.Println("error listing resources:", err)
		http.Error(w, "Error listing resources", http.StatusInternalServerError)
		return
	}

	_ = resourceview.ResourcesPage(&resourceview.ResourcesPageProps{
		Ctx:            ctx,
		ShowArchived:   uv.IsArchived,
		Resources:      resources,
		ResourcesCount: count,
		Sort:           sort,
		Page:           uv.Page,
		PageSize:       uv.PageSize,
	}).Render(w)
}

type resourceHomePageURLVals struct {
	Sort     string
	Page     int
	PageSize int

	IsArchived bool
}

func (uv *resourceHomePageURLVals) normalise() {
	if uv.Page == 0 {
		uv.Page = 1
	}
	if uv.PageSize == 0 {
		uv.PageSize = 50
	}
}

type resourcePageURLVals struct {
	ShowArchivedSchedules bool
}

func (h *ResourceHandler) ResourcePage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	canUserEdit := ctx.User.Permissions.UserAdmin.Access

	var resourcePageUV resourcePageURLVals
	if err := appurl.Unmarshal(r.URL.Query(), &resourcePageUV); err != nil {
		log.Println("error decoding resource page query params:", err)
		http.Error(w, "Error decoding query params", http.StatusBadRequest)
		return
	}

	resourceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
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

	currentMetrics, err := h.resourceService.GetResourceServiceSchedules(
		r.Context(),
		resourceID,
		resourcePageUV.ShowArchivedSchedules)
	if err != nil {
		log.Println("error fetching resource metrics summary:", err)
		http.Error(w, "Error fetching resource metrics summary", http.StatusInternalServerError)
		return
	}

	serviceHistoryQuery := model.GetServicesQuery{
		Page:       1,
		PageSize:   100,
		ResourceIn: []string{resource.Reference},
	}

	services, serviceCount, err := h.servicesService.ListServices(r.Context(), serviceHistoryQuery)
	if err != nil {
		log.Println("error listing resource services:", err)
		http.Error(w, "Error listing resource services", http.StatusInternalServerError)
		return
	}

	lifetimeTotals, err := h.resourceService.GetServiceMetricLifetimeTotals(
		r.Context(),
		resourceID,
	)
	if err != nil {
		log.Println("error fetching lifetime totals:", err)
		http.Error(w, "Error fetching lifetime totals", http.StatusInternalServerError)
		return
	}

	for i, service := range services {
		if canUserEdit {
			services[i].GalleryURL = h.galleryService.GenerateEditTempURL(service.GalleryID, true)
		} else {
			services[i].GalleryURL = h.galleryService.GenerateTempURL(service.GalleryID, canUserEdit)
		}
	}

	canManageSchedules := canUserEdit && !resource.IsArchived

	_ = resourceview.ResourcePage(&resourceview.ResourcePageProps{
		Ctx:                   ctx,
		Resource:              *resource,
		Services:              services,
		CurrentMetrics:        currentMetrics,
		LifetimeTotals:        lifetimeTotals,
		ServiceCount:          serviceCount,
		Sort:                  serviceHistoryQuery.Sort,
		Page:                  serviceHistoryQuery.Page,
		PageSize:              serviceHistoryQuery.PageSize,
		CanManageSchedules:    canManageSchedules,
		ShowArchivedSchedules: resourcePageUV.ShowArchivedSchedules,
	}).Render(w)
}

func (h *ResourceHandler) AddResourcePage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	type urlVals struct {
		Type       string
		Reference  string
		IsArchived bool
	}

	var uv urlVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding query params", http.StatusInternalServerError)
		return
	}

	teams, _, err := h.teamService.List(r.Context(), model.ListTeamsQuery{
		Page: 1, PageSize: 1000,
	})
	if err != nil {
		log.Println("error fetching teams:", err)
		http.Error(w, "Error fetching teams", http.StatusInternalServerError)
		return
	}

	_ = resourceview.AddResourcePage(&resourceview.AddResourcePageProps{
		Ctx:    ctx,
		Values: r.URL.Query(),
		Teams:  teams,
	}).Render(w)
}

func (h *ResourceHandler) AddResource(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	type urlVals struct {
		Type       string
		Reference  string
		IsArchived bool
	}

	var uv urlVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		log.Println("error decoding query params:", err)
		http.Error(w, "Error decoding query params", http.StatusBadRequest)
		return
	}

	teams, _, err := h.teamService.List(r.Context(), model.ListTeamsQuery{
		Page: 1, PageSize: 1000,
	})
	if err != nil {
		log.Println("error fetching teams:", err)
		http.Error(w, "Error fetching teams", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println("error parsing form:", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var fd addResourceFormData
	if err := appurl.Unmarshal(r.Form, &fd); err != nil {
		log.Println("error decoding form:", err)
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd.normalise()

	validationErrors := fd.validate()

	if len(validationErrors) > 0 {

		_ = resourceview.AddResourcePage(&resourceview.AddResourcePageProps{
			Ctx:              ctx,
			Values:           r.Form,
			Teams:            teams,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	err = h.resourceService.CreateResource(
		r.Context(),
		model.NewResource{
			Type:                   fd.Type,
			Reference:              fd.Reference,
			ServiceOwnershipTeamID: fd.ServiceOwnershipTeamID,
		})
	if err != nil {
		log.Println("error creating resource:", err)
		http.Error(w, "Error creating resource", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/resources", http.StatusSeeOther)
}

type addResourceFormData struct {
	Type                   string
	Reference              string
	ServiceOwnershipTeamID *int
}

func (fd *addResourceFormData) normalise() {
	fd.Type = strings.ToUpper(strings.TrimSpace(fd.Type))
	fd.Reference = strings.TrimSpace(fd.Reference)
}

func (fd *addResourceFormData) validate() validate.ValidationErrors {
	var ve validate.ValidationErrors = make(map[string][]string)

	if fd.Type == "" {
		ve.Add("Type", "is required")
	}

	if fd.Reference == "" {
		ve.Add("Reference", "is required")
	}

	return ve
}

func (h *ResourceHandler) EditResourcePage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	resourceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
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

	teams, _, err := h.teamService.List(r.Context(), model.ListTeamsQuery{
		Page: 1, PageSize: 1000,
	})
	if err != nil {
		log.Println("error fetching teams:", err)
		http.Error(w, "Error fetching teams", http.StatusInternalServerError)
		return
	}

	_ = resourceview.EditResourcePage(&resourceview.EditResourcePageProps{
		Ctx:      ctx,
		Resource: *resource,
		Teams:    teams,
	}).Render(w)
}

func (h *ResourceHandler) EditResource(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	resourceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid resource id", http.StatusBadRequest)
		return
	}

	teams, _, err := h.teamService.List(r.Context(), model.ListTeamsQuery{
		Page: 1, PageSize: 1000,
	})
	if err != nil {
		log.Println("error fetching teams:", err)
		http.Error(w, "Error fetching teams", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println("error parsing form:", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var fd editResourceFormData
	if err := appurl.Unmarshal(r.Form, &fd); err != nil {
		log.Println("error decoding form:", err)
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd.normalise()

	validationErrors := fd.validate()
	if len(validationErrors) > 0 {
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

		_ = resourceview.EditResourcePage(&resourceview.EditResourcePageProps{
			Ctx:              ctx,
			Resource:         *resource,
			Teams:            teams,
			Values:           r.Form,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	serviceValidationErrors, err := h.resourceService.UpdateResource(r.Context(), resourceID, model.ResourceUpdate{
		Type:                   fd.Type,
		Reference:              fd.Reference,
		IsArchived:             fd.IsArchived,
		ServiceOwnershipTeamID: fd.ServiceOwnershipTeamID,
	})
	if err != nil {
		log.Println("error updating resource:", err)
		http.Error(w, "Error updating resource", http.StatusInternalServerError)
		return
	}
	if len(serviceValidationErrors) > 0 {
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

		_ = resourceview.EditResourcePage(&resourceview.EditResourcePageProps{
			Ctx:              ctx,
			Resource:         *resource,
			Teams:            teams,
			Values:           r.Form,
			ValidationErrors: serviceValidationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/resources/%d", resourceID), http.StatusSeeOther)
}

type editResourceFormData struct {
	Type                   string
	Reference              string
	IsArchived             bool
	ServiceOwnershipTeamID *int
}

func (fd *editResourceFormData) normalise() {
	fd.Type = strings.ToUpper(strings.TrimSpace(fd.Type))
	fd.Reference = strings.TrimSpace(fd.Reference)
}

func (fd *editResourceFormData) validate() validate.ValidationErrors {
	var ve validate.ValidationErrors = make(map[string][]string)

	if fd.Type == "" {
		ve.Add("Type", "is required")
	}

	if fd.Reference == "" {
		ve.Add("Reference", "is required")
	}

	return ve
}

func (h *ResourceHandler) AddResourceServicePage(
	w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	resourceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
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
	if resource == nil || resource.IsArchived {
		http.Error(w, "Resource not available", http.StatusNotFound)
		return
	}

	_ = resourceview.AddResourceServicePage(
		&resourceview.AddResourceServicePageProps{
			Ctx:      ctx,
			Resource: *resource,
		}).Render(w)
}

func (h *ResourceHandler) AddResourceService(
	w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	resourceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
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

	var fd addResourceServiceFormData
	if err := appurl.Unmarshal(r.Form, &fd); err != nil {
		log.Println("error decoding form:", err)
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd.normalise()

	validationErrors := fd.validate()

	if len(validationErrors) > 0 {
		_ = resourceview.AddResourceServicePage(
			&resourceview.AddResourceServicePageProps{
				Ctx:              ctx,
				Values:           r.Form,
				ValidationErrors: validationErrors,
				IsSubmission:     true,
				Resource:         *resource,
			}).Render(w)
		return
	}

	fd.ResourceID = resourceID

	serviceID, err := h.resourceService.CreateResourceService(
		r.Context(),
		model.NewResourceService{
			ResourceID: fd.ResourceID,
			Notes:      fd.Notes,
		}, ctx.User.UserID)
	if err != nil {
		log.Println("error creating resource service:", err)
		http.Error(w, "Error creating resource service", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/services/%d", serviceID), http.StatusSeeOther)
}

type addResourceServiceFormData struct {
	ResourceID int
	Notes      string
}

func (fd *addResourceServiceFormData) normalise() {
}

func (fd *addResourceServiceFormData) validate() validate.ValidationErrors {
	var ve validate.ValidationErrors = make(map[string][]string)

	return ve
}

func (h *ResourceHandler) AddResourceMetricRecordPage(
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
	if resource == nil || resource.IsArchived {
		http.Error(w, "Resource not available", http.StatusNotFound)
		return
	}

	metrics, err := h.resourceService.GetResourceServiceMetrics(r.Context(), resourceID)
	if err != nil {
		log.Println("error fetching resource service metrics:", err)
		http.Error(w, "Error fetching resource service metrics", http.StatusInternalServerError)
		return
	}

	_ = resourceview.AddResourceMetricRecordPage(&resourceview.AddResourceMetricRecordPageProps{
		Ctx:            ctx,
		Values:         r.URL.Query(),
		Resource:       *resource,
		ServiceMetrics: metrics,
	}).Render(w)
}

func (h *ResourceHandler) AddResourceMetricRecord(w http.ResponseWriter, r *http.Request) {
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

	metrics, err := h.resourceService.GetResourceServiceMetrics(r.Context(), resourceID)
	if err != nil {
		log.Println("error fetching resource service metrics:", err)
		http.Error(w, "Error fetching resource service metrics", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println("error parsing form:", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var fd addResourceMetricRecordFormData
	if err := appurl.Unmarshal(r.Form, &fd); err != nil {
		log.Println("error decoding form:", err)
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd.normalise()

	validationErrors := fd.validate()

	if len(validationErrors) > 0 {
		_ = resourceview.AddResourceMetricRecordPage(
			&resourceview.AddResourceMetricRecordPageProps{
				Ctx:              ctx,
				Values:           r.Form,
				ValidationErrors: validationErrors,
				IsSubmission:     true,
				Resource:         *resource,
				ServiceMetrics:   metrics,
			}).Render(w)
		return
	}

	err = h.resourceService.CreateResourceMetricRecord(
		r.Context(),
		model.NewResourceServiceMetricRecord{
			ResourceID:              resourceID,
			ResourceServiceMetricID: fd.ServiceMetricID,
			Value:                   fd.Value,
		})
	if err != nil {
		log.Println("error creating resource recording:", err)
		http.Error(w, "Error creating resource recording", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/resources/%d", resourceID), http.StatusSeeOther)
}

type addResourceMetricRecordFormData struct {
	ServiceMetricID int
	Value           decimal.Decimal
}

func (fd *addResourceMetricRecordFormData) normalise() {
}

func (fd *addResourceMetricRecordFormData) validate() validate.ValidationErrors {
	var ve validate.ValidationErrors = make(map[string][]string)

	if fd.ServiceMetricID == 0 {
		ve.Add("ServiceMetricID", "must be selected")
	}
	if !fd.Value.GreaterThan(decimal.Zero) {
		ve.Add("Value", "should be positive")
	}

	return ve
}
