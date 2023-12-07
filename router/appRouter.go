package router

import (
	"net/http"
	"operationalcore/components"
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

	// Form router
	r.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
		_ = views.Form().Render(w)
	})

	r.HandleFunc("/options", func(w http.ResponseWriter, r *http.Request) {
		// Get the search value from the request
		searchValue := r.FormValue("search-value")

		var options = []components.Option{
			{
				Value: searchValue + searchValue,
				Label: searchValue + searchValue,
			},
			{
				Value: searchValue + searchValue + searchValue,
				Label: searchValue + searchValue + searchValue,
			},
			{
				Value: searchValue + searchValue + searchValue + searchValue,
				Label: searchValue + searchValue + searchValue + searchValue,
			},
		}

		_ = components.MultiSelectOptions(options, "search-select-dropdown", "").Render(w)
	})

	// module routers
	addContactsRouter(r)

	// static assets file server
	staticFS := http.FileServer(static.Assets)
	r.PathPrefix("/").Handler(staticFS)

	return r
}
