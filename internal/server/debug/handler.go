package debug

import (
	"github.com/albertwidi/go_project_example/internal/pkg/router"
	"github.com/albertwidi/go_project_example/internal/server"
	"github.com/albertwidi/go_project_example/internal/server/debug/user"
)

// Handlers of debug server
type Handlers struct {
	user *user.Handler
}

func (s *Server) handler(debugHandlers Handlers) *router.Router {
	r := router.New(nil)
	r.Use(server.Metrics)
	r.Get("/someting", debugHandlers.user.BypassLogin)

	return r
}
