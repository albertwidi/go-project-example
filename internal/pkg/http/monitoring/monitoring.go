// this package provide delagation of http.ResponseWriter for monitoring purpose

package monitoring

import "net/http"

// Delegator interface for delegating http.ResponseWriter
type Delegator interface {
	http.ResponseWriter
	Status() int
	Written() int64
}

// responseWriterDelagator is a http.ResponseWriter delegator
type responseWriterDelegator struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

// Status return the status of WriteHeader
func (rwdg *responseWriterDelegator) Status() int {
	return rwdg.status
}

func (rwdg *responseWriterDelegator) Written() int64 {
	return rwdg.written
}

// WriteHeader via http.ResponseWriter
func (rwdg *responseWriterDelegator) WriteHeader(code int) {
	if rwdg.wroteHeader {
		return
	}
	rwdg.status = code
	rwdg.wroteHeader = true
	rwdg.ResponseWriter.WriteHeader(code)
}

// Write the byte via http.ResponseWriter
func (rwdg *responseWriterDelegator) Write(b []byte) (int, error) {
	if !rwdg.wroteHeader {
		rwdg.WriteHeader(http.StatusOK)
	}

	n, err := rwdg.ResponseWriter.Write(b)
	rwdg.written += int64(n)
	return n, err
}

// NewResponseWriterDelegator is a delegator for http.ResponseWriter
// with some additional function for monitoring purpose
func NewResponseWriterDelegator(w http.ResponseWriter) Delegator {
	d := responseWriterDelegator{
		ResponseWriter: w,
	}
	return &d
}
