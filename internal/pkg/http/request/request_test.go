package request

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestNewRequest(t *testing.T) {
	cases := []struct {
		name         string
		url          string
		method       string
		urlQuery     map[string]string
		postForm     map[string]string
		header       map[string]string
		noVersion    bool
		expectHeader map[string]string
		expectBody   string
		expectError  bool
	}{
		{
			name:   "simple request test",
			url:    "http://service2.cluster.local",
			method: http.MethodGet,
			header: map[string]string{
				"route-version-select": "service1.cluster.local|0.1.2+beta,service2.cluster.local|0.2.0",
			},
			expectHeader: map[string]string{
				"routes-version-select": "service1.cluster.local|0.1.2+beta,service2.cluster.local|0.2.0",
				"route-version-select":  "service1.cluster.local|0.1.2+beta",
				"version-select":        "0.2.0",
			},
			expectError: false,
		},
		{
			name:   "simple request with no version",
			url:    "http://service2.cluster.local",
			method: http.MethodGet,
			header: map[string]string{
				"route-version-select": "service1.cluster.local|0.1.2+beta,service2.cluster.local|0.2.0",
			},
			noVersion: true,
			expectHeader: map[string]string{
				"routes-version-select": "",
				"route-version-select":  "",
				"version-select":        "",
			},
			expectError: false,
		},
		{
			name:   "url with rest path",
			url:    "http://service2.cluster.local/v1/testing/10",
			method: http.MethodGet,
			header: map[string]string{
				"route-version-select": "service1.cluster.local|0.1.2+beta,service2.cluster.local|0.2.0",
			},
			expectHeader: map[string]string{
				"routes-version-select": "service1.cluster.local|0.1.2+beta,service2.cluster.local|0.2.0",
				"route-version-select":  "service1.cluster.local|0.1.2+beta",
				"version-select":        "0.2.0",
			},
			expectError: false,
		},
		{
			name:   "url with port and rest path",
			url:    "http://service2.cluster.local:9000/v1/testing/10",
			method: http.MethodGet,
			header: map[string]string{
				"route-version-select": "service1.cluster.local|0.1.2+beta,service2.cluster.local|0.2.0",
			},
			noVersion: false,
			expectHeader: map[string]string{
				"routes-version-select": "service1.cluster.local|0.1.2+beta,service2.cluster.local|0.2.0",
				"route-version-select":  "service1.cluster.local|0.1.2+beta",
				"version-select":        "0.2.0",
			},
			expectError: false,
		},
		{
			name:   "post request with no version",
			url:    "http://service2.cluster.local",
			method: http.MethodPost,
			postForm: map[string]string{
				"key1": "val1",
				"key2": "val2",
			},
			header: map[string]string{
				"route-version-select": "service1.cluster.local|0.1.2+beta,service2.cluster.local| 0.2.0",
			},
			noVersion: true,
			expectHeader: map[string]string{
				"routes-version-select": "",
				"route-version-select":  "",
				"version-select":        "",
			},
			expectBody:  "key1=val1&key2=val2",
			expectError: false,
		},
		{
			name:   "post request with body postform",
			url:    "http://service2.cluster.local",
			method: http.MethodPost,
			postForm: map[string]string{
				"key1": "val1",
				"key2": "val2",
			},
			header: map[string]string{
				"route-version-select": "service1.cluster.local|0.1.2+beta,service2.cluster.local| 0.2.0",
			},
			expectHeader: map[string]string{
				"routes-version-select": "service1.cluster.local|0.1.2+beta,service2.cluster.local| 0.2.0",
				"route-version-select":  "service1.cluster.local|0.1.2+beta",
				"version-select":        "0.2.0",
			},
			expectBody:  "key1=val1&key2=val2",
			expectError: false,
		},
	}

	for _, c := range cases {
		t.Log(c.name)
		ctx := context.WithValue(context.Background(), &RoutingContext, c.header["route-version-select"])
		g := New(ctx).
			Method(c.method).
			URL(c.url)
		if c.noVersion {
			g.NoVersionHeader()
		}

		switch c.method {
		case http.MethodGet:
		case http.MethodPost:
			kv := make([]string, 0)
			for k, v := range c.postForm {
				kv = append(kv, k, v)
			}
			g.PostForm(kv...)
		}

		req, err := g.Compile()
		if err != nil {
			t.Error(err)
			return
		}

		// check error
		if !c.expectError && err != nil {
			t.Errorf("expecting no error but got %v", err)
			return
		}

		// check header and generated header
		for k, v := range c.expectHeader {
			h := req.Header.Get(k)
			if h != v {
				t.Errorf("expecting %s but got %s for %s header", v, h, k)
				return
			}
		}

		// check url
		if c.url != req.URL.String() {
			t.Errorf("expecting url %s but got %s", c.url, req.URL.String())
			return
		}

		// check url query

		// check body
		if c.expectBody != "" {
			out, err := ioutil.ReadAll(req.Body)
			if err != nil {
				t.Error(err)
				return
			}

			if string(out) != c.expectBody {
				t.Errorf("expecting body %s but got %s", c.expectBody, string(out))
				return
			}
		}
	}
}
