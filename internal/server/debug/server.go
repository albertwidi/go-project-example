// debug server is for debug/development purpose only
// to provide some debug functionality
// for example login, file server, etc

package debug

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
