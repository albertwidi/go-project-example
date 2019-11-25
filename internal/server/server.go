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
	Run() error
}

// Addresses of server
type Addresses struct {
	Main  string
	Admin string
	Debug string
}

// Server configuration
type Server struct {
	Address string

	// prometheus vector object for metrics
	countervec      *prometheus.CounterVec
	durationhist    *prometheus.HistogramVec
	requestsizehist *prometheus.HistogramVec
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
