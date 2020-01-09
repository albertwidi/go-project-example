package server

import (
	"context"
	"net"
	"net/http"

	"github.com/albertwidi/go-project-example/internal/pkg/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type adminServer struct {
	address    string
	httpServer *http.Server
	listener   net.Listener
}

func (s *Server) newAdminServer(address string) (*adminServer, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	adm := adminServer{
		address:    address,
		listener:   listener,
		httpServer: &http.Server{},
	}
	return &adm, nil
}

// Run admin server
func (adm *adminServer) Run(middlewares ...router.MiddlewareFunc) error {
	r := router.New(adm.address, nil)
	r.Use(middlewares...)
	adm.registerHandler(r)
	adm.httpServer.Handler = r
	return adm.httpServer.Serve(adm.listener)
}

// Shutdown admin server
func (adm *adminServer) Shutdown(ctx context.Context) error {
	return adm.httpServer.Shutdown(ctx)
}

func (adm *adminServer) registerHandler(r *router.Router) {
	r.Handle("/metrics", promhttp.Handler())
}
