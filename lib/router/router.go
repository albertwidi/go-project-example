package router

import (
	"log"
	"net/http"
	"strings"

	requestcontext "github.com/albertwidi/go_project_example/lib/context"
	"github.com/albertwidi/go_project_example/lib/http/misc"
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
	router *mux.Router
	route  []*mux.Route
	mw     []MiddlewareFunc
	// options
	options Options
}

// Options of router
type Options struct {
	Debug bool
}

// New router
func New(options Options) *Router {
	r := Router{
		router:  mux.NewRouter(),
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
func (r *Router) Use(handler MiddlewareFunc) {
	r.mw = append(r.mw, handler)
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
	route.Methods("GET")
	route.Path(path)
	r.handleRoute(route, "get", handler)
}

// Head function
func (r *Router) Head(path string, handler HandlerFunc) {
	route := r.router.NewRoute()
	route.Methods("HEAD")
	route.Path(path)
	r.handleRoute(route, "head", handler)
}

// Post function
func (r *Router) Post(path string, handler HandlerFunc) {
	route := r.router.NewRoute()
	route.Methods("POST")
	route.Path(path)
	r.handleRoute(route, "post", handler)
}

// Patch function
func (r *Router) Patch(path string, handler HandlerFunc) {
	route := r.router.NewRoute()
	route.Methods("PATCH")
	route.Path(path)
	r.handleRoute(route, "patch", handler)
}

// Delete function
func (r *Router) Delete(path string, handler HandlerFunc) {
	route := r.router.NewRoute()
	route.Methods("DELETE")
	route.Path(path)
	r.handleRoute(route, "delete", handler)
}

// Options function
func (r *Router) Options(path string, handler HandlerFunc) {
	route := r.router.NewRoute()
	route.Methods("OPTIONS")
	route.Path(path)
	r.handleRoute(route, "options", handler)
}

// Handle request with pure http handler
func (r *Router) Handle(path string, handler http.Handler) {
	r.router.Handle(path, handler)
}

// PathPrefix implementation of mux router
func (r *Router) PathPrefix(tpl string) *mux.Route {
	return r.router.PathPrefix(tpl)
}

// handleRoute function
func (r *Router) handleRoute(route *mux.Route, method string, handler HandlerFunc) {
	r.route = append(r.route, route)
	pathTemplate, err := route.GetPathTemplate()
	if err != nil {
		log.Fatal(err)
	}

	handlerFunc := func(writer http.ResponseWriter, request *http.Request) {
		requestContext := requestcontext.New(requestcontext.Constructor{
			HTTPResponseWriter: writer,
			HTTPRequest:        request,
			Path:               pathTemplate,
			Method:             misc.SanitizeMethod(method),
		})

		h := handler
		for i := range r.mw {
			h = r.mw[len(r.mw)-1-i](h)
		}
		h(requestContext)
	}

	route.HandlerFunc(handlerFunc)
}

// ServeHTTP function
func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	r.router.ServeHTTP(writer, request)
}
