package router

import (
	"app/internal/handler"
	"app/internal/service"
	"app/pkg/apphmac"
	"net/http"
)

func addServiceRoutes(
	mux *http.ServeMux,
	galleryService service.GalleryService,
	resourceService service.ResourceService,
	servicesService service.ServicesService,
	commentService service.CommentService,
	teamService service.TeamService,
	appHMAC apphmac.AppHMAC,
) {

	servicesHandler := handler.NewServiceHandler(
		resourceService,
		servicesService,
		galleryService,
		commentService,
		teamService,
		appHMAC,
	)

	mux.HandleFunc("GET /services/metrics", servicesHandler.ServiceMetricsPage)

	mux.HandleFunc("GET /services/metrics/add", servicesHandler.AddResourceServiceMetricPage)
	mux.HandleFunc("POST /services/metrics/add", servicesHandler.AddResourceServiceMetric)
	mux.HandleFunc("GET /services/metrics/{id}/edit", servicesHandler.EditResourceServiceMetricPage)
	mux.HandleFunc("POST /services/metrics/{id}/edit", servicesHandler.EditResourceServiceMetric)
	mux.HandleFunc("PUT /services/{serviceID}/resource/{resourceID}/{action}/update", servicesHandler.UpdateResourceService)

	mux.HandleFunc("GET /services/schedules", servicesHandler.ServiceSchedulesPage)

	mux.HandleFunc("GET /services/schedules/add", servicesHandler.AddServiceSchedulePage)
	mux.HandleFunc("POST /services/schedules/add", servicesHandler.AddServiceSchedule)
	mux.HandleFunc("GET /services/schedules/{id}/edit", servicesHandler.EditServiceSchedulePage)
	mux.HandleFunc("POST /services/schedules/{id}/edit", servicesHandler.EditServiceSchedule)

	mux.HandleFunc("GET /services/resource/{id}/schedules/add", servicesHandler.AddResourceServiceSchedulePage)
	mux.HandleFunc("POST /services/resource/{id}/schedules/add", servicesHandler.AssignServiceSchedule)
	mux.HandleFunc("POST /services/resource/{id}/schedules/{scheduleID}/unassign", servicesHandler.UnassignServiceSchedule)

	mux.HandleFunc("GET /services", servicesHandler.ResourceServicingPage)
	mux.HandleFunc("GET /services/all", servicesHandler.ServicesPage)

	mux.HandleFunc("GET /services/{serviceID}", servicesHandler.ResourceServicePage)
	mux.HandleFunc("POST /services/{serviceID}", servicesHandler.UpdateResourceServiceNotes)

}
