package context

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/albertwidi/go-project-example/internal/pkg/http/response"
)

// RequestContext struct
type RequestContext struct {
	httpResponseWriter http.ResponseWriter
	httpRequest        *http.Request
	address            string
	path               string
	method             string
}

// Constructor of context
type Constructor struct {
	HTTPResponseWriter http.ResponseWriter
	HTTPRequest        *http.Request
	Address            string
	Path               string
	Method             string
}

// New context
func New(constructor Constructor) *RequestContext {
	rc := RequestContext{
		httpResponseWriter: constructor.HTTPResponseWriter,
		httpRequest:        constructor.HTTPRequest,
		address:            constructor.Address,
		path:               constructor.Path,
		method:             constructor.Method,
	}
	return &rc
}

// Address return the address where request arrived to
func (rc *RequestContext) Address() string {
	return rc.address
}

// Request return http request from request context
func (rc *RequestContext) Request() *http.Request {
	return rc.httpRequest
}

// RequestHeader return http.Request.Header
func (rc *RequestContext) RequestHeader() http.Header {
	return rc.httpRequest.Header
}

// RequestHandler return handler name of the request
func (rc *RequestContext) RequestHandler() string {
	return rc.path
}

// Context return the http.Request.Context
func (rc *RequestContext) Context() context.Context {
	return rc.httpRequest.Context()
}

// ResponseWriter return http response writer from request context
func (rc *RequestContext) ResponseWriter() http.ResponseWriter {
	return rc.httpResponseWriter
}

// JSON to create a json response via http response lib
func (rc *RequestContext) JSON() *response.JSONResponse {
	j := response.JSON(rc.httpResponseWriter)
	return j
}

// DecodeJSON from request body
func (rc *RequestContext) DecodeJSON(out interface{}) error {
	in, err := ioutil.ReadAll(rc.httpRequest.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(in, out); err != nil {
		return err
	}
	return nil
}
