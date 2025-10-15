package router

import (
	"app/internal/handler"
	"app/internal/service"
	"app/pkg/apphmac"
	"net/http"
)

func addAndonRoutes(
	mux *http.ServeMux,
	andonService service.AndonService,
	andonIssueService service.AndonIssueService,
	commentService service.CommentService,
	fileService service.FileService,
	galleryService service.GalleryService,
	teamService service.TeamService,
	appHMAC apphmac.AppHMAC,
) {
	andonHandler := handler.NewAndonHandler(
		andonService,
		andonIssueService,
		commentService,
		fileService,
		galleryService,
		teamService,
		appHMAC,
	)

	mux.HandleFunc("GET /andons", andonHandler.HomePage)
	mux.HandleFunc("GET /andons/all", andonHandler.AllAndonsPage)

	mux.HandleFunc("GET /andons/add", andonHandler.AddPage)
	mux.HandleFunc("POST /andons/add", andonHandler.Add)

	mux.HandleFunc("GET /andons/{andonID}", andonHandler.AndonPage)

	mux.HandleFunc("POST /andons/{andonID}/{action}/update", andonHandler.UpdateAndon)

}
