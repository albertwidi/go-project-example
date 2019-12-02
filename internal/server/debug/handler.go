package debug

import (
	"github.com/albertwidi/go_project_example/internal/pkg/router"
	"github.com/albertwidi/go_project_example/internal/server/debug/user"
)

// Handlers of debug server
type Handlers struct {
	user *user.Handler
}

func (s *Server) registerHandlers(r *router.Router) {
	r.Get("/something", s.handlers.user.BypassLogin)
}
