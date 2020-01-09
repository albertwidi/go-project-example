package server

import (
	"context"
	"time"

	requestctx "github.com/albertwidi/go-project-example/internal/pkg/context"
	httpmisc "github.com/albertwidi/go-project-example/internal/pkg/http/misc"
	httpmonitoring "github.com/albertwidi/go-project-example/internal/pkg/http/monitoring"
	"github.com/albertwidi/go-project-example/internal/pkg/router"
	"github.com/prometheus/client_golang/prometheus"
)

// Runner server interface
type Runner interface {
	// Run and pass mandatory middlewares when running each server
	Run(middlewares ...router.MiddlewareFunc) error
	Shutdown(ctx context.Context) error
}

// Server configuration
type Server struct {
	runners []Runner
	errChan chan error

	// prometheus vector object for metrics
	countervec      *prometheus.CounterVec
	durationhist    *prometheus.HistogramVec
	requestsizehist *prometheus.HistogramVec
}

// Run the server
func (s *Server) Run() chan error {
	for _, r := range s.runners {
		go func(r Runner) {
			if err := r.Run(s.Metrics); err != nil {
				s.errChan <- err
			}
		}(r)
	}
	return s.errChan
}

// Shutdown the server
// TODO: check whether one of the runners return error when shutting down
func (s *Server) Shutdown(ctx context.Context) error {
	// send nil error to get the server out of the loop
	s.errChan <- nil

	for _, r := range s.runners {
		r.Shutdown(ctx)
	}
	return nil
}

// New server
func New(adminServerAddress string, runners ...Runner) (*Server, error) {
	// initialize monitoring metrics
	countervec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "a counter for requests to the wrapped handler.",
		},
		[]string{"address", "code", "method", "path"},
	)
	err := prometheus.DefaultRegisterer.Register(countervec)
	if err != nil {
		return nil, err
	}
	durationhist := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration",
			Help: "a histogram of request latencies.",
		},
		[]string{"address", "code", "method", "path"},
	)
	err = prometheus.DefaultRegisterer.Register(durationhist)
	if err != nil {
		return nil, err
	}
	requestsizehist := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_size",
			Help: "a histrogram of request size",
		},
		[]string{"address", "code", "method", "path"},
	)
	err = prometheus.DefaultRegisterer.Register(requestsizehist)
	if err != nil {
		return nil, err
	}
	s := Server{
		errChan:         make(chan error, 1),
		countervec:      countervec,
		durationhist:    durationhist,
		requestsizehist: requestsizehist,
		runners:         runners,
	}

	if s.runners == nil {
		s.runners = []Runner{}
	}
	adm, err := s.newAdminServer(adminServerAddress)
	if err != nil {
		return nil, err
	}
	s.runners = append(s.runners, adm)
	return &s, nil
}

// Metrics is a middleware for metrics monitoring
func (s *Server) Metrics(next router.HandlerFunc) router.HandlerFunc {
	return func(rctx *requestctx.RequestContext) error {
		now := time.Now()
		err := next(rctx)
		httpStatus := 0
		duration := time.Since(now).Seconds()
		requestSize := httpmisc.ComputeApproximateRequestSize(rctx.Request())
		requestMethod := rctx.Request().Method
		handlerName := rctx.RequestHandler()
		address := rctx.Address()

		w := rctx.ResponseWriter()
		d, ok := w.(httpmonitoring.Delegator)
		if ok {
			httpStatus = d.Status()
		}
		s.countervec.WithLabelValues(address, httpmisc.SanitizeCode(httpStatus), requestMethod, handlerName).Inc()
		s.durationhist.WithLabelValues(address, httpmisc.SanitizeCode(httpStatus), requestMethod, handlerName).Observe(duration)
		s.requestsizehist.WithLabelValues(address, httpmisc.SanitizeCode(httpStatus), requestMethod, handlerName).Observe(float64(requestSize))
		return err
	}
}

// Trace is a middleware for tracing using opentelemetry
func Trace(next router.HandlerFunc) router.HandlerFunc {
	return func(rctx *requestctx.RequestContext) error {
		return nil
	}
}
