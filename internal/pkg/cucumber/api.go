package cucumber

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/gherkin"
)

// APIFeature to test api
type APIFeature struct {
	Options APIFeatureOptions
	logger  *log.Logger
	client  *http.Client
	// request
	requestHeader http.Header
	requestBody   io.Reader
	// response
	// we might want to change this to io.Reader
	// to have a more flexible response control
	responseBody   []byte
	responseCode   int
	responseHeader http.Header
}

// APIFeatureOptions contain options of api feature
type APIFeatureOptions struct {
	// EndpointMapping provides mapping feature to the gherkin value
	// for example, we don't have to always state the full endpoint
	EndpointsMapping map[string]string
}

func (api *APIFeature) reset() {
	api.client = &http.Client{}
	api.requestHeader = http.Header{}
	api.requestBody = nil
	api.responseBody = nil
	api.responseCode = 0
	api.responseHeader = http.Header{}
}

func (api *APIFeature) setRequestHeader(header string) error {
	headers := strings.Split(header, ":")
	if len(headers) != 2 {
		return fmt.Errorf("invalid headers length, got %v", headers)
	}
	api.requestHeader = http.Header{}

	api.requestHeader.Add(strings.TrimSpace(headers[0]), strings.TrimSpace(headers[1]))
	return nil
}

func (api *APIFeature) setRequestBody(body *gherkin.DocString) error {
	if body.Content == "" {
		return errors.New("body is empty")
	}
	if api.requestHeader.Get("Content-Type") == "application/json" {
		rawjson := json.RawMessage(body.Content)
		out, err := rawjson.MarshalJSON()
		if err != nil {
			return fmt.Errorf("SetRequestbody: %w", err)
		}
		api.requestBody = bytes.NewBuffer(out)
	} else {
		api.requestBody = bytes.NewBufferString(body.Content)
	}
	return nil
}

func (api *APIFeature) isendRequestToWithPath(method, service, path string) error {
	e, ok := api.Options.EndpointsMapping[service]
	if !ok {
		return fmt.Errorf("api_feature: endpoint with name %s not exist", service)
	}
	endpoint := mergeEndpointAndPath(e, path)
	return api.isendRequestTo(method, endpoint)
}

func mergeEndpointAndPath(endpoint, path string) string {
	// return nothing if endpoint is not specified
	if endpoint == "" {
		return ""
	}
	// normalize the endpoint and path to make it less error prone
	// for example if endpoint is http://127.0.0.1/ and path is /v1/book/detail
	// it will become http://127.0.0.1//v1/book/detail, normalize to make things consistent
	if endpoint[len(endpoint)-1:] == "/" {
		endpoint = endpoint[:len(endpoint)-1]
	}
	if len(path) > 0 {
		if path[:1] != "/" {
			path = "/" + path[0:]
		}
	}
	return endpoint + path
}

func (api *APIFeature) isendRequestTo(method, endpoint string) error {
	req, err := http.NewRequest(method, endpoint, api.requestBody)
	if err != nil {
		return err
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}

	return api.setHTTPResponseComponents(resp)
}

func (api *APIFeature) setHTTPResponseComponents(resp *http.Response) error {
	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	api.responseCode = resp.StatusCode
	api.responseBody = out
	api.responseHeader = resp.Header
	return nil
}

func (api *APIFeature) theResponseCodeShouldBe(code int) error {
	if code != api.responseCode {
		return fmt.Errorf("api: expect response code %d but got %d", code, api.responseCode)
	}
	return nil
}

func (api *APIFeature) theResponseHeaderShouldBe(headerName, headerValue string) error {
	v := api.responseHeader.Get(headerName)
	if v != headerValue {
		return fmt.Errorf("header: expect value %s but got %s instead from header %s", headerValue, v, headerName)
	}
	return nil
}

func (api *APIFeature) theResponseShouldMatch(body *gherkin.DocString) error {
	if string(api.responseBody) != body.Content {
		return fmt.Errorf("expect output of %s but got %s", string(api.responseBody), body.Content)
	}
	return nil
}

func (api *APIFeature) theResponseShouldMatchJSON(body *gherkin.DocString) error {
	var expected, actual interface{}

	// re-encode expected response
	if err := json.Unmarshal([]byte(body.Content), &expected); err != nil {
		return err
	}

	// re-encode actual response too
	if err := json.Unmarshal(api.responseBody, &actual); err != nil {
		return err
	}

	// the matching may be adapted per different requirements.
	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expected, actual)
	}
	return nil
}

// BeforeRegister for invoking things needed before registering feature to cucumber
func (api *APIFeature) BeforeRegister() (err error) {
	return
}

// SetLogger to set default logger
func (api *APIFeature) SetLogger(logger *log.Logger) {
	api.logger = logger
}

// FeatureContext for api
func (api *APIFeature) FeatureContext(s *godog.Suite) {
	s.BeforeScenario(func(v interface{}) {
		api.reset()
	})
	// Given steps
	s.Step(`^set request header "([^"]*)"$`, api.setRequestHeader)
	s.Step(`^set request body:$`, api.setRequestBody)
	// When steps
	s.Step(`^I send "(GET|POST|PUT|PATCH|DELETE)" request to "([^"]*)"$`, api.isendRequestTo)
	s.Step(`^I send "(GET|POST|PUT|PATCH|DELETE)" request to "([^"]*) service" with path "([^"]*)"$`, api.isendRequestToWithPath)
	// Then steps
	s.Step(`^dump the response body to file "([^"]*)"$`, api.theResponseCodeShouldBe)
	s.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	s.Step(`^the response header "([^"]*)" should be "([^"]*)"$`, api.theResponseHeaderShouldBe)
	s.Step(`^the response should match json:$`, api.theResponseShouldMatchJSON)
	s.Step(`^the response should match body:$`, api.theResponseShouldMatch)
}
