package middleware

import (
	"net/http"
	"strings"
)

var publicRoutes = []string{
	"/static/",
	"/auth/",
}

func isPublicRouteRequest(r *http.Request) bool {
	for _, route := range publicRoutes {
		if strings.HasPrefix(r.URL.Path, route) {
			return true
		}
	}
	return false
}
