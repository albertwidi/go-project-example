package router

import (
	"log"
	"net/http"
	"strings"

	requestcontext "github.com/albertwidi/go-project-example/internal/pkg/context"
	"github.com/albertwidi/go-project-example/internal/pkg/http/misc"
	"github.com/albertwidi/go-project-example/internal/pkg/http/monitoring"
	"github.com/gorilla/mux"
)

// HandlerFunc to handle http request
type HandlerFunc func(ctx *requestcontext.RequestContext) error

// MiddlewareFunc to handle middleware chaining
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// Next function
func (mw MiddlewareFunc) Next(handler HandlerFunc) HandlerFunc {
	return mw(handler)
}

// ChainedMiddlewares type
type ChainedMiddlewares struct {
	middlewares []MiddlewareFunc
	r           *Router
}

// NewChainedMiddleware function
func NewChainedMiddleware(r *Router, middlewares ...MiddlewareFunc) *ChainedMiddlewares {
	return &ChainedMiddlewares{
		middlewares: middlewares,
		r:           r,
	}
}

// Then function
func (cmw *ChainedMiddlewares) Then(method, path string, handler HandlerFunc) {
	for _, middleware := range cmw.middlewares {
		handler = middleware.Next(handler)
	}
	cmw.r.HandleFunc(method, path, handler)
}

// Get right from router
func (cmw *ChainedMiddlewares) Get(path string, handler HandlerFunc) {
	for _, middleware := range cmw.middlewares {
		handler = middleware.Next(handler)
	}
	cmw.r.Get(path, handler)
}

// Post right from router
func (cmw *ChainedMiddlewares) Post(path string, handler HandlerFunc) {
	for _, middleware := range cmw.middlewares {
		handler = middleware.Next(handler)
	}
	cmw.r.Post(path, handler)
}

// Delete right from router
func (cmw *ChainedMiddlewares) Delete(path string, handler HandlerFunc) {
	for _, middleware := range cmw.middlewares {
		handler = middleware.Next(handler)
	}
	cmw.r.Delete(path, handler)
}

// Patch right from router
func (cmw *ChainedMiddlewares) Patch(path string, handler HandlerFunc) {
	for _, middleware := range cmw.middlewares {
		handler = middleware.Next(handler)
	}
	cmw.r.Patch(path, handler)
}

// Head right from router
func (cmw *ChainedMiddlewares) Head(path string, handler HandlerFunc) {
	for _, middleware := range cmw.middlewares {
		handler = middleware.Next(handler)
	}
	cmw.r.Head(path, handler)
}

// Options right from router
func (cmw *ChainedMiddlewares) Options(path string, handler HandlerFunc) {
	for _, middleware := range cmw.middlewares {
		handler = middleware.Next(handler)
	}
	cmw.r.Options(path, handler)
}

// Router struct
type Router struct {
	router  *mux.Router
	route   []*mux.Route
	mw      []MiddlewareFunc
	address string
	// options
	options *Options
}

// Options of router
type Options struct {
	Debug bool
}

// New router
func New(address string, options *Options) *Router {
	if options == nil {
		options = &Options{
			Debug: false,
		}
	}

	r := Router{
		router:  mux.NewRouter(),
		address: address,
		options: options,
	}
	return &r
}

// Vars return map of query parameters
func (r *Router) Vars(req *http.Request) map[string]string {
	return mux.Vars(req)
}

// Routes return list of route in the router
func (r *Router) Routes() []*mux.Route {
	return r.route
}

// Use middleware
func (r *Router) Use(middlewares ...MiddlewareFunc) {
	for _, m := range middlewares {
		r.mw = append(r.mw, m)
	}
}

// HandleFunc function
func (r *Router) HandleFunc(method, path string, handler HandlerFunc) {
	method = strings.ToUpper(method)
	route := r.router.NewRoute()
	route.Methods(method)
	route.Path(path)
	r.handleRoute(route, method, handler)
}

// Get function
func (r *Router) Get(path string, handler HandlerFunc) {
	route := r.router.NewRoute()
	method := http.MethodGet
	route.Methods(method)
	route.Path(path)
	r.handleRoute(route, method, handler)
}

// Head function
func (r *Router) Head(path string, handler HandlerFunc) {
	route := r.router.NewRoute()
	method := http.MethodHead
	route.Methods(method)
	route.Path(path)
	r.handleRoute(route, method, handler)
}

// Post function
func (r *Router) Post(path string, handler HandlerFunc) {
	route := r.router.NewRoute()
	method := http.MethodPost
	route.Methods(method)
	route.Path(path)
	r.handleRoute(route, method, handler)
}

// Patch function
func (r *Router) Patch(path string, handler HandlerFunc) {
	route := r.router.NewRoute()
	method := http.MethodPatch
	route.Methods(method)
	route.Path(path)
	r.handleRoute(route, method, handler)
}

// Delete function
func (r *Router) Delete(path string, handler HandlerFunc) {
	route := r.router.NewRoute()
	method := http.MethodDelete
	route.Methods(method)
	route.Path(path)
	r.handleRoute(route, method, handler)
}

// Options function
func (r *Router) Options(path string, handler HandlerFunc) {
	route := r.router.NewRoute()
	method := http.MethodOptions
	route.Methods(method)
	route.Path(path)
	r.handleRoute(route, method, handler)
}

// Handle request with pure http handler
func (r *Router) Handle(path string, handler http.Handler) {
	route := r.router.NewRoute()
	route.Path(path)
	r.handleRoute(route, path, handler)
}

// PathPrefix implementation of mux router
// this function is not using custom handleRoute functions
// metrics/diagnostics will not exported from this method
func (r *Router) PathPrefix(tpl string) *mux.Route {
	return r.router.PathPrefix(tpl)
}

// handleRoute function
// handleRoute always use http.ResponseWriter delegator
func (r *Router) handleRoute(route *mux.Route, method string, hn interface{}) {
	r.route = append(r.route, route)
	pathTemplate, err := route.GetPathTemplate()
	if err != nil {
		log.Fatal(err)
	}

	switch v := hn.(type) {
	case HandlerFunc:
		handlerFunc := func(writer http.ResponseWriter, request *http.Request) {
			// always use http.ResponseWriter delegator for monitoring purpose
			delegator := monitoring.NewResponseWriterDelegator(writer)
			requestContext := requestcontext.New(requestcontext.Constructor{
				HTTPResponseWriter: delegator,
				HTTPRequest:        request,
				Address:            r.address,
				Path:               pathTemplate,
				Method:             misc.SanitizeMethod(method),
			})

			h := v
			for i := range r.mw {
				h = r.mw[len(r.mw)-1-i](h)
			}
			h(requestContext)
		}
		route.HandlerFunc(handlerFunc)

	case http.HandlerFunc:
		f := func(w http.ResponseWriter, req *http.Request) {
			delegator := monitoring.NewResponseWriterDelegator(w)
			v(delegator, req)
		}
		route.HandlerFunc(f)

	case http.Handler:
		route.Handler(v)
	}
}

// ServeHTTP function
func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	r.router.ServeHTTP(writer, request)
}
