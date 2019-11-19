package server

import "context"

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
