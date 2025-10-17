package main

import (
	"net/http"

	"github.com/DTSL/golang-libraries/httpjson"
	"github.com/DTSL/golang-libraries/httpsib"
	"github.com/DTSL/golang-libraries/httptracing/httptracinggorilla"
	"github.com/gorilla/mux"
	"go.uber.org/dig"
)

// NewRouterParams defines NewRouter parameters.
type newRouterParams struct {
	dig.In

	HTTPTestHandler *httpTestHandler
}

// endpoint defines router endpoint.
type endpoint struct {
	Path   string
	Routes []func(route *mux.Route)
}

func newRouter(params newRouterParams) http.Handler {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(httpNotFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(httpMethodNotAllowed)
	for _, e := range []endpoint{
		// Define API endpoints here
		{
			Path: "/api",
			Routes: []func(*mux.Route){
				func(r *mux.Route) {
					r.Methods(http.MethodGet).Handler(params.HTTPTestHandler)
				},
			},
		},
	} {
		for _, fn := range e.Routes {
			fn(router.NewRoute().Path(e.Path))
		}
	}
	return httpsib.WrapHandler(
		httptracinggorilla.WrapRouter(router),
		httpsib.OTELMetrics(appname, version),
	)
}

func httpNotFound(w http.ResponseWriter, req *http.Request) {
	_ = httpjson.WriteResponse(req.Context(), w, http.StatusNotFound, nil)
}

func httpMethodNotAllowed(w http.ResponseWriter, req *http.Request) {
	_ = httpjson.WriteResponse(req.Context(), w, http.StatusMethodNotAllowed, nil)
}
