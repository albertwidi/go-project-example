package debug

// go:generate swagger

import (
	"github.com/albertwidi/go-project-example/internal/pkg/router"
	"github.com/albertwidi/go-project-example/internal/server/debug/user"
)

// Handlers of debug server
type Handlers struct {
	user *user.Handler
}

func (s *Server) registerHandlers(r *router.Router) {
	// swagger:route GET /user/login/bypass bypass user login
	// Bypassing user login
	// This will bypass user login
	// Only ued in development
	//	Consumes:
	//	- application/json
	//	Produces:
	//	- application/json
	//	Schemes: http
	r.Get("/user/login/bypass", s.handlers.user.BypassLogin)
}
