package request

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var (
	// flag to activate request url and forwarded header
	_requestVersionMatching = true

	// RoutingContext is a key name for version selection routing
	RoutingContext = "REQUEST_ROUTING_HEADER"
)

// Request wrap the http request
type Request struct {
	method           string
	url              string
	query            string
	header           http.Header
	additionalHeader []string

	body      io.Reader
	vBody     interface{}
	bodyJSON  bool
	noVersion bool

	ctx context.Context
	r   *http.Request
}

// New http request wrapper
func New(ctx context.Context) *Request {
	r := Request{ctx: ctx}
	return &r
}

// Headers function to set request header
func (r *Request) Headers(header http.Header) *Request {
	r.header = header
	return r
}

// NoVersionHeader function to control the version selection header generation
func (r *Request) NoVersionHeader() *Request {
	r.noVersion = true
	return r
}

// Method function set the request method
func (r *Request) Method(method string) *Request {
	r.method = method
	return r
}

// URL function set the request url
func (r *Request) URL(url string) *Request {
	r.url = url
	return r
}

// Get function for building get request
func (r *Request) Get(url string) *Request {
	r.method = http.MethodGet
	r.url = url
	return r
}

// Query for creating url query
func (r *Request) Query(kv ...string) *Request {
	data := url.Values{}
	for idx := range kv {
		if idx > 0 {
			idx++
			if idx == len(kv)-1 {
				break
			}
		}
		data.Add(kv[idx], kv[idx+1])
	}
	r.query = data.Encode()
	return r
}

// Post function for building post request
func (r *Request) Post(url string) *Request {
	r.method = http.MethodPost
	return r
}

// PostForm set a url values for a postform body in a request
func (r *Request) PostForm(kv ...string) *Request {
	data := url.Values{}
	for idx := range kv {
		if idx > 0 {
			idx++
			if idx == len(kv)-1 {
				break
			}
		}
		data.Add(kv[idx], kv[idx+1])
	}
	// expected to create a body, and not append the postform to save allocations
	r.body = strings.NewReader(data.Encode())
	return r
}

// Put function for building put request
func (r *Request) Put(url string) *Request {
	r.method = http.MethodPut
	return r
}

// Body of the request
func (r *Request) Body(body io.Reader) *Request {
	r.body = body
	return r
}

// BodyJSON indicate that request body is a json data
func (r *Request) BodyJSON(body interface{}) *Request {
	r.vBody = body
	r.bodyJSON = true
	r.additionalHeader = []string{"Content-Type", "application/json"}
	return r
}

// Compile the http request
// version selection header spesification:
// 1. routes-version-select
//		propagated header for selecting version in matching url
//		example: svc1.cluster.local|0.1.0, svc2.cluster.local|0.2.0
// 2. route-version-select
//		propagated header for selecting version in matching url
//		different from routes-version-select, matching url will not propagated
//		example: svc1.cluster.local|0.1.0
// 3. version-select
//		propageted header when an url is match, only contain version value
//		example: 0.2.0
func (r *Request) Compile() (*http.Request, error) {
	u, err := url.Parse(r.url)
	if err != nil {
		return nil, err
	}

	finalURL := u.String()
	if r.query != "" {
		finalURL = finalURL + "?" + r.query
	}
	req, err := http.NewRequestWithContext(r.ctx, r.method, finalURL, r.body)
	if err != nil {
		return nil, err
	}
	req.Header = r.header

	// flag the version matching
	// this logic might moved to infrastructure instead here
	// we can easily disable this without any side-effect
	if _requestVersionMatching && !r.noVersion {
		rvcHeader := ""
		vheader, vroutes := getRoutingHeader(r.ctx)
		for k, v := range vroutes {
			if k == u.Hostname() {
				// add header for version select because url is matching
				// so the infrastructure able to route via header
				req.Header.Set("version-select", v)
			} else {
				// add header for route-version-select
				// all unmatched value needs to be appended
				// avoid fmt.Sprintf because it will allocate
				rvcHeader += k + "|" + v + ","
			}
		}
		// trim most right "," token from the header
		rvcHeader = rvcHeader[:len(rvcHeader)-1]
		req.Header.Set("route-version-select", rvcHeader)
		req.Header.Set("routes-version-select", vheader)
	}
	return req, nil
}

// getRoutingHeader return the routings header and decoded version of routings header
// routing header means to what version that the request is directed to
// with this special header, we want to enable A/B testing and specific version selection in request
// for example reqeust to service1.cluster.local with header:
// [route-version-select] service2.cluster.local|0.1.2,service3.cluster.local|0.2
// will inform service1 to contact service 2 with version 0.1.2 and service with version 0.2
func getRoutingHeader(ctx context.Context) (string, map[string]string) {
	v := ctx.Value(&RoutingContext)
	header, ok := v.(string)
	// make sure header is valid and not empty
	if !ok || header == "" {
		return "", nil
	}
	// multiple specification in one header/context is allowed
	// and the version selection is using smver for the specification
	// spec: https://semver.org/
	// for example, it is allowed to specifiy 1.0.0-beta+exp.sha.5114f85
	// as the semver is not using ':' token in the spesificaion, we can use
	// strings.Split to divide the value
	routings := strings.Split(header, ",")
	vroutings := make(map[string]string)
	for idx, r := range routings {
		selections := strings.Split(r, "|")
		// continue, because the spefication is broken
		if len(selections) < 2 {
			continue
		}
		if idx > 1 {
			idx += 2
		}
		vroutings[strings.TrimSpace(selections[0])] = strings.TrimSpace(selections[1])
	}
	return header, vroutings
}
