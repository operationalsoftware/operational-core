package router

import (
	"app/internal/handler"
	"app/internal/service"
	"app/pkg/apphmac"
	"net/http"
)

func addResourceRoutes(
	mux *http.ServeMux,
	galleryService service.GalleryService,
	resourceService service.ResourceService,
	servicesService service.ServicesService,
	commentService service.CommentService,
	teamService service.TeamService,
	appHMAC apphmac.AppHMAC,
) {

	resourceHandler := handler.NewResourceHandler(
		resourceService,
		servicesService,
		galleryService,
		commentService,
		teamService,
		appHMAC,
	)

	mux.HandleFunc("GET /resources", resourceHandler.ResourcesPage)

	mux.HandleFunc("GET /resources/add", resourceHandler.AddResourcePage)
	mux.HandleFunc("POST /resources/add", resourceHandler.AddResource)

	mux.HandleFunc("GET /resources/{id}", resourceHandler.ResourcePage)
	mux.HandleFunc("GET /resources/{id}/edit", resourceHandler.EditResourcePage)
	mux.HandleFunc("POST /resources/{id}/edit", resourceHandler.EditResource)

	mux.HandleFunc("GET /resources/{id}/services/new", resourceHandler.AddResourceServicePage)
	mux.HandleFunc("POST /resources/{id}/services/new", resourceHandler.AddResourceService)

	mux.HandleFunc("GET /resources/{id}/usage/add", resourceHandler.AddResourceUsagePage)
	mux.HandleFunc("POST /resources/{id}/usage/add", resourceHandler.AddResourceUsage)
}
