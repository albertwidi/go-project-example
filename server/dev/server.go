// dev server is for development purpose only
// only active when PROJECTENV=development
// this is to provide some testing functionality
// for example login, file server, etc

package dev

// Server struct
type Server struct {
	address string
}

// New server
func New(address string) (*Server, error) {
	s := Server{}
	return &s, nil
}

// Run dev server
func (s *Server) Run() error {
	return nil
}
