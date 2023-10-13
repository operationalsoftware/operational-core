package router

import (
	"net/http"
	"operationalcore/static"
	"operationalcore/views"

	"github.com/gorilla/mux"
)

func AppRouter() *mux.Router {
	r := mux.NewRouter()

	// homepage router
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = views.Index().Render(w)
	})

	// module routers
	addContactsRouter(r)

	// static assets file server
	staticFS := http.FileServer(static.Assets)
	r.PathPrefix("/").Handler(staticFS)

	return r
}
