package response

import (
	"encoding/json"
	"net/http"

	"github.com/albertwidi/go-project-example/internal/xerrors"
)

// Status of json response
type Status string

// list of status for http response
const (
	StatusOK            Status = "OK"
	StatusRetry         Status = "RETRY"
	StatusBadRequest    Status = "BAD_REQUEST"
	StatusNotFound      Status = "NOT_FOUND"
	StatusUnauthorized  Status = "UNAUTHORIZED"
	StatusInternalError Status = "INTERNAL_ERROR"
)

// JSONResponse struct for http json response
type JSONResponse struct {
	writer        http.ResponseWriter
	xerr          *xerrors.Errors
	headerWritten bool

	// response part
	ResponseStatus Status             `json:"status"`
	ResponseData   interface{}        `json:"data"`
	ResponseRetry  *JSONRetryResponse `json:"retry,omitempty"`
	ResponseError  *JSONError         `json:"error,omitempty"`
}

// JSONRetryResponse struct for retry field in json response
type JSONRetryResponse struct {
	RetryMin int `json:"retry_min"`
	RetryMax int `json:"retry_max"`
}

// JSONError struct for error field in json response
type JSONError struct {
	Title   string   `json:"title"`
	Message string   `json:"message"`
	Detail  string   `json:"detail"`
	Errors  []string `json:"errors"`
}

// JSON create a new JSON response
func JSON(w http.ResponseWriter) *JSONResponse {
	resp := JSONResponse{
		writer: w,
	}
	return &resp
}

// SetHeader used to set header in http.ResponseWriter of JSONResponse
func (jresp *JSONResponse) SetHeader(key, value string) {
	jresp.writer.Header().Set(key, value)
}

// Data for set data to json response
func (jresp *JSONResponse) Data(data interface{}) *JSONResponse {
	jresp.ResponseData = data
	return jresp
}

// Error set error to json response
// only use error when the type of error is *xerrors.Error
func (jresp *JSONResponse) Error(err error, errResp *JSONError) *JSONResponse {
	xerr, ok := err.(*xerrors.Errors)
	if !ok {
		return jresp
	}
	jresp.xerr = xerr
	jresp.ResponseError = errResp
	return jresp
}

// WriteHeader set the header of JSONResponse writer
func (jresp *JSONResponse) WriteHeader(statusCode int) *JSONResponse {
	if jresp.headerWritten {
		return jresp
	}
	jresp.writer.WriteHeader(statusCode)
	jresp.headerWritten = true
	return jresp
}

// Write json response
func (jresp *JSONResponse) Write() (int, error) {
	jresp.writer.Header().Set("Content-Type", "application/json")
	// process the error internals
	if jresp.xerr != nil {
		kind := jresp.xerr.Kind()
		switch kind {
		case xerrors.KindOK:
			jresp.ResponseStatus = StatusOK
			jresp.WriteHeader(http.StatusOK)

		case xerrors.KindNotFound:
			jresp.ResponseStatus = StatusNotFound
			jresp.WriteHeader(http.StatusNotFound)

		case xerrors.KindBadRequest:
			jresp.ResponseStatus = StatusBadRequest
			jresp.writer.WriteHeader(http.StatusBadRequest)

		case xerrors.KindUnauthorized:
			jresp.ResponseStatus = StatusUnauthorized
			jresp.WriteHeader(http.StatusUnauthorized)

		case xerrors.KindInternalError:
			jresp.ResponseStatus = StatusInternalError
			jresp.WriteHeader(http.StatusInternalServerError)
		}
	}

	out, err := json.Marshal(jresp)
	if err != nil {
		return 0, err
	}
	return jresp.writer.Write(out)
}
