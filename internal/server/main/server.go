package mainserver

import "context"

// Server struct
type Server struct {
	address string
}

// New http server
func New(address string) {

}

// Run the http server
func (s *Server) Run() error {
	return nil
}

// Shutdown the main server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}
