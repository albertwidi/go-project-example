package debug

import (
	"github.com/albertwidi/go_project_example/internal/pkg/router"
	"github.com/albertwidi/go_project_example/internal/server/debug/user"
)

// Handlers of debug server
type Handlers struct {
	user *user.Handler
}

func (s *Server) handler(debugHandlers Handlers, middlewares ...router.MiddlewareFunc) *router.Router {
	r := router.New(s.address, nil)
	r.Use(middlewares...)
	r.Get("/someting", debugHandlers.user.BypassLogin)

	return r
}
