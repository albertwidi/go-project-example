// debug server is for debug/development purpose only
// to provide some debug functionality
// for example login, file server, etc

package debug

import (
	"context"
	"net"
	"net/http"

	"github.com/albertwidi/go-project-example/debug/user"
	"github.com/albertwidi/go-project-example/internal/pkg/router"
	userhandler "github.com/albertwidi/go-project-example/internal/server/debug/user"
)

// Server struct
type Server struct {
	address    string
	httpServer *http.Server
	listener   net.Listener
	// handlers
	handlers Handlers
}

// Usecases of debug server
type Usecases struct {
	user *user.DebugUsecase
}

// New server
func New(address string, usecases Usecases) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	// init all handlers
	userHandler := userhandler.New(usecases.user)
	handlers := Handlers{
		user: userHandler,
	}
	s := Server{
		address:    address,
		listener:   listener,
		httpServer: &http.Server{},
		handlers:   handlers,
	}
	return &s, nil
}

// Run debug server
func (s *Server) Run(middlewares ...router.MiddlewareFunc) error {
	// initiate httpserver handler
	r := router.New(s.address, nil)
	r.Use(middlewares...)
	s.registerHandlers(r)
	s.httpServer.Handler = r
	return s.httpServer.Serve(s.listener)
}

// Shutdown debug server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
