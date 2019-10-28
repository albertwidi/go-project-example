package http

import (
	"context"
	"net"
	"net/http"

	"github.com/albertwidi/go_project_example/api"
	"github.com/albertwidi/go_project_example/api/http/order"
	"github.com/albertwidi/go_project_example/api/http/payment"
	"github.com/albertwidi/go_project_example/api/http/user"
	"github.com/gorilla/mux"
)

type Server struct {
	server         *http.Server
	UserService    api.UserService
	OrderService   api.OrderService
	PaymentService api.PaymentService
}

func (s *Server) Serve(lis net.Listener) error {
	s.server = &http.Server{}

	// init all handler
	user.Init(s.UserService)
	order.Init(s.OrderService)
	payment.Init(s.PaymentService)
	// import all route into server handler
	s.server.Handler = handler()

	return s.server.Serve(lis)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func handler() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/user/hello", user.Hello).Methods("GET")
	return r
}
