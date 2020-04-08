package cucumber

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cucumber/godog/gherkin"
)

func TestReset(t *testing.T) {
	api := APIFeature{
		responseBody: []byte("haloha"),
	}
	api.reset()
	if api.responseBody != nil {
		t.Error("object reset failed")
		return
	}
}

func TestMergeEndpointAndPath(t *testing.T) {
	cases := []struct {
		name     string
		endpoint string
		path     string
		result   string
	}{
		{
			name:     "one",
			endpoint: "http://127.0.0.1/",
			path:     "/v1/book/detail",
			result:   "http://127.0.0.1/v1/book/detail",
		},
		{
			name:     "two",
			endpoint: "http://127.0.0.1",
			path:     "v1/book/detail",
			result:   "http://127.0.0.1/v1/book/detail",
		},
		{
			name:     "three",
			endpoint: "http://127.0.0.1/",
			path:     "v1/book/detail",
			result:   "http://127.0.0.1/v1/book/detail",
		},
		{
			name:     "four",
			endpoint: "http://127.0.0.1",
			path:     "/v1/book/detail",
			result:   "http://127.0.0.1/v1/book/detail",
		},
	}

	for _, c := range cases {
		result := mergeEndpointAndPath(c.endpoint, c.path)
		if result != c.result {
			t.Errorf("%s: expecting result %s but got %s", c.name, c.result, result)
			return
		}
	}
}

func TestApiFeatureHTTPHeader(t *testing.T) {
	cases := []struct {
		header      string
		expectKey   string
		expectValue string
	}{
		{
			header:      "Content-Type: application/json",
			expectKey:   "Content-Type",
			expectValue: "application/json",
		},
		{
			header:      "Content-Type: application/x-www-form-urlencoded",
			expectKey:   "Content-Type",
			expectValue: "application/x-www-form-urlencoded",
		},
		{
			header:      "Authorization: bearer jfhk20281037201js",
			expectKey:   "Authorization",
			expectValue: "bearer jfhk20281037201js",
		},
	}

	for _, c := range cases {
		api := APIFeature{}
		if err := api.setRequestHeader(c.header); err != nil {
			t.Error(err)
			return
		}

		value := api.requestHeader.Get(c.expectKey)
		if value != c.expectValue {
			t.Errorf("%s: expecting %s header value for %s but got %s", t.Name(), c.expectValue, c.expectKey, value)
			return
		}
	}
}

// TestIsSendRequest functio to test the godog tester
// the function to test itself should be correct and tested
// this flow is tested in this function:
// 1. Create a http request
// 2. Check the http response code
// 3. Check the http response body
// 4. Check the http header
//
// NOTES:
// - Invalid JSON request/response can produce unhelpful error, need to check whether json is valid/not
func TestIsSendRequestTo(t *testing.T) {
	cases := []struct {
		// request
		requestEndpoint string
		requestMethod   string
		requestHeader   map[string]string
		requestBody     string
		// response
		responseCode   int
		responseHeader map[string]string
		responseJSON   bool
		responseBody   string
	}{
		// a simple request with simple request body
		{
			// request
			requestEndpoint: "/v1/something",
			requestMethod:   "POST",
			requestHeader: map[string]string{
				"Content-Type": "application/json",
			},
			requestBody: `
			{
				"key": "value"
			}
			`,
			// response
			responseCode: http.StatusOK,
			responseHeader: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
			responseJSON: true,
			responseBody: `
			{
				"key": "value"
			}
			`,
		},
		// a simple request with more complex response body
		{
			// request
			requestEndpoint: "/v1/something",
			requestMethod:   "POST",
			requestHeader: map[string]string{
				"Content-Type": "application/json",
			},
			requestBody: `
			{
				"key": "value"
			}
			`,
			// response
			responseCode: http.StatusOK,
			responseHeader: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
			responseJSON: true,
			responseBody: `
			{
				"key": {
					"key": {
						"key": "value"
					},
					"key2": {
						"key": "value"
					},
					"key3": 10,
					"key4": 10.1
				},
				"key2": {
					"key": {
						"key": "value",
						"key2": false
					}
				}
			}
			`,
		},
	}

	for _, c := range cases {
		api := APIFeature{}
		api.reset()
		handler := func(w http.ResponseWriter, r *http.Request) {
			var (
				out []byte
				err error
			)

			if strings.ToUpper(r.Method) != strings.ToUpper(c.requestMethod) {
				err = fmt.Errorf("expect method %s but got %s", c.requestMethod, r.Method)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			for k, v := range c.responseHeader {
				w.Header().Add(k, v)
			}

			out = []byte(c.responseBody)
			if c.responseJSON {
				raw := json.RawMessage([]byte(c.responseBody))
				out, err = raw.MarshalJSON()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
					return
				}
			}
			w.WriteHeader(c.responseCode)
			w.Write(out)
		}

		var header = http.Header{}
		for k, v := range c.requestHeader {
			header.Add(k, v)
		}

		var requestBody io.Reader
		// check whether the body is json or url encoded
		// if json, then we need to set the body from json.RawMessage
		if header.Get("Content-Type") == "application/json" {
			raw := json.RawMessage([]byte(c.requestBody))
			b, err := raw.MarshalJSON()
			if err != nil {
				t.Error(err)
				return
			}
			requestBody = bytes.NewBuffer(b)
		} else if header.Get("Content-Type") == "application/x-www-form-urlencoded" {
			requestBody = bytes.NewBufferString(c.requestBody)
		}

		// create a test request and recorder
		req := httptest.NewRequest(c.requestMethod, c.requestEndpoint, requestBody)
		req.Header = header
		w := httptest.NewRecorder()
		handler(w, req)

		// get the http response from the http.ResponseWriter recorder
		resp := w.Result()
		// setHttpResponseComponents is a function to set components of http response to APIFeature object
		// this is important, because we analyze the response based on this function
		if err := api.setHTTPResponseComponents(resp); err != nil {
			t.Error(err)
			return
		}

		// test whether the response code is correct
		if err := api.theResponseCodeShouldBe(resp.StatusCode); err != nil {
			t.Error(err)
			return
		}

		// check the response output
		if c.responseJSON {
			raw := json.RawMessage([]byte(c.responseBody))
			jsonb, err := raw.MarshalJSON()
			if err != nil {
				t.Error(err)
				return
			}
			if err := api.theResponseShouldMatch(&gherkin.DocString{Content: string(jsonb)}); err != nil {
				t.Error(err)
				return
			}
		} else {
			if err := api.theResponseShouldMatch(&gherkin.DocString{Content: c.responseBody}); err != nil {
				t.Error(err)
				return
			}
		}

		// check the response header matching
		for k, v := range c.responseHeader {
			if err := api.theResponseHeaderShouldBe(k, v); err != nil {
				t.Error(err)
				return
			}
		}
	}
}
