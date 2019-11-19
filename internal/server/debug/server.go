// debug server is for debug/development purpose only
// to provide some debug functionality
// for example login, file server, etc

package debug

import (
	"context"
	"net"
	"net/http"

	"github.com/albertwidi/go_project_example/debug/user"
	userhandler "github.com/albertwidi/go_project_example/internal/server/debug/user"
)

// Server struct
type Server struct {
	address    string
	httpServer *http.Server
	listener   net.Listener
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

	userHandlers := userhandler.New(usecases.user)
	handlers := Handlers{
		user: userHandlers,
	}

	s := Server{
		address:    address,
		listener:   listener,
		httpServer: &http.Server{},
	}
	// attatch all handlers to http server
	s.httpServer.Handler = s.handler(handlers)
	return &s, nil
}

// Run dev server
func (s *Server) Run() error {
	return s.httpServer.Serve(s.listener)
}

// Shutdown debug server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
