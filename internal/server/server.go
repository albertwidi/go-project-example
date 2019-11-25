package server

import (
	"context"
	requestctx "github.com/albertwidi/go_project_example/internal/pkg/context"
	"github.com/albertwidi/go_project_example/internal/pkg/router"
)

// Addresses of server
type Addresses struct {
	Main  string
	Admin string
	Debug string
}

// Server configuration
type Server struct {
	Address string
}

// Run the server
func (s *Server) run() error {
	return nil
}

// Shutdown the server
func (s *Server) shutdown(ctx context.Context) error {
	return nil
}

// Usecases for the server
type Usecases struct {
}

// New server
func New() {

}

// Metrics is a middleware for metrics monitoring
func Metrics(next router.HandlerFunc) router.HandlerFunc {
	return func(rctx *requestctx.RequestContext) error {
		return nil
	}
}
