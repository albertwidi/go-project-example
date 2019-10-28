package grpc

import (
	"net"

	"github.com/albertwidi/go_project_example/api"
	"google.golang.org/grpc"
)

type Server struct {
	server  *grpc.Server
	user    api.UserService
	order   api.OrderService
	payment api.PaymentService
}

// Start grpc service
func (s *Server) Serve(lis net.Listener) error {
	s.server = grpc.NewServer()
	return s.server.Serve(lis)
}

func (s *Server) Shutdown() error {
	s.server.GracefulStop()
	return nil
}
