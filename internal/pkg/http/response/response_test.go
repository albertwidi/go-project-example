package response_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/albertwidi/go_project_example/internal/pkg/http/response"
	"github.com/albertwidi/go_project_example/internal/xerrors"
)

func kindToStatusCode(err *xerrors.Errors) int {
	switch err.Kind() {
	case xerrors.KindOK:
		return http.StatusOK
	case xerrors.KindBadRequest:
		return http.StatusBadRequest
	case xerrors.KindNotFound:
		return http.StatusNotFound
	case xerrors.KindUnauthorized:
		return http.StatusUnauthorized
	case xerrors.KindInternalError:
		return http.StatusInternalServerError
	}
	return 0
}

func TestWrite(t *testing.T) {
	cases := []struct {
		Name       string
		Headers    map[string]string
		HTTPStatus int
		XErrors    error
	}{
		{
			Name:       "1",
			HTTPStatus: http.StatusOK,
		},
		{
			Name:    "2",
			XErrors: xerrors.New("bad request", xerrors.KindBadRequest),
		},
	}

	for _, c := range cases {
		t.Logf("test number: %s", c.Name)
		handler := func(w http.ResponseWriter, r *http.Request) {
			jsonresp := response.JSON(w)
			if c.XErrors == nil {
				jsonresp.WriteHeader(c.HTTPStatus)
			} else {
				jsonresp.Error(c.XErrors, nil)
			}
			jsonresp.Write()
		}

		req := httptest.NewRequest("GET", "http://example.com", nil)
		w := httptest.NewRecorder()
		handler(w, req)

		resp := w.Result()
		statusCode := c.HTTPStatus
		if c.XErrors != nil {
			// always expect *xerrors.Errors
			statusCode = kindToStatusCode(c.XErrors.(*xerrors.Errors))
		}
		// check status code
		t.Logf("status code: %d", statusCode)
		if statusCode != resp.StatusCode {
			t.Errorf("invalid http status, expect %d but got %d", c.HTTPStatus, resp.StatusCode)
			return
		}
		// check header
		for key, val := range c.Headers {
			hval := resp.Header.Get(key)
			if hval != val {
				t.Errorf("invalid header value, expect %s but got %s", val, hval)
				return
			}
		}
	}
}
