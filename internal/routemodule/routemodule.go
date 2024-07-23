package routemodule

import "net/http"

type RouteModule interface {
	AddRoutes(r *http.ServeMux, prefix string)
	GetPrefix() string
}

type PrefixedRouteModule struct {
	Prefix string
}

func (b *PrefixedRouteModule) GetPrefix() string {
	return b.Prefix
}
