package server

import (
	"context"
	"time"

	requestctx "github.com/albertwidi/go_project_example/internal/pkg/context"
	httpmisc "github.com/albertwidi/go_project_example/internal/pkg/http/misc"
	httpmonitoring "github.com/albertwidi/go_project_example/internal/pkg/http/monitoring"
	"github.com/albertwidi/go_project_example/internal/pkg/router"
	"github.com/prometheus/client_golang/prometheus"
)

// Runner server interface
type Runner interface {
	// Run and pass mandatory middlewares when running each server
	Run(middlewares ...router.MiddlewareFunc) error
	Shutdown(ctx context.Context) error
}

// Addresses of server
type Addresses struct {
	Main  string
	Admin string
	Debug string
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
func (s *Server) Run() error {
	for _, r := range s.runners {
		go func(r Runner) {
			if err := r.Run(s.Metrics); err != nil {
				s.errChan <- err
			}
		}(r)
	}

	for {
		err := <-s.errChan
		if err != nil {
			return err
		}
	}
}

// Shutdown the server
// TODO: check whether one of the runners return error when shutting down
func (s *Server) Shutdown(ctx context.Context) error {
	for _, r := range s.runners {
		r.Shutdown(ctx)
	}
	return nil
}

// New server
func New(runner ...Runner) (*Server, error) {
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
	}
	return &s, nil
}

// Metrics is a middleware for metrics monitoring
func (s *Server) Metrics(next router.HandlerFunc) router.HandlerFunc {
	return func(rctx *requestctx.RequestContext) error {
		now := time.Now()

		d := httpmonitoring.NewResponseWriterDelegator(rctx.ResponseWriter())
		rctx.SetResponseWriter(d)
		err := next(rctx)

		duration := time.Since(now).Seconds()
		requestSize := httpmisc.ComputeApproximateRequestSize(rctx.Request())
		requestMethod := rctx.Request().Method
		handlerName := rctx.RequestHandler()
		address := rctx.Address()

		s.countervec.WithLabelValues(address, httpmisc.SanitizeCode(d.Status()), requestMethod, handlerName).Inc()
		s.durationhist.WithLabelValues(address, httpmisc.SanitizeCode(d.Status()), requestMethod, handlerName).Observe(duration)
		s.requestsizehist.WithLabelValues(address, httpmisc.SanitizeCode(d.Status()), requestMethod, handlerName).Observe(float64(requestSize))
		return err
	}
}

// Trace is a middleware for tracing using opentelemetry
func Trace(next router.HandlerFunc) router.HandlerFunc {
	return func(rctx *requestctx.RequestContext) error {
		return nil
	}
}
