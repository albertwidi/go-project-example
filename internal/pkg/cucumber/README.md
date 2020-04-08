# Cucumber

Test Framework Using BDD(Behavior Driven Development) Pattern. To know more about cucumber, please go to this [site](https://cucumber.io/)

This framework is based on `cucumber/godog` library which leverage `gherkin` grammar for testing.

## Limitiations

This testing framework is intended to test the API for integration test. It is possible to include the integration for database, but currently not supported.

## Installing Godog

To install godog, you need `go`. If `go` is already installed, then please execute this script:

```shell
go get github.com/cucumber/godog/cmd/godog
```

## Run The Test

To run the test within this directory, run:

```shell
go test -v *.go
```

## Supported Grammar

(Given, When, Then - Gherkin):

- `Given` some context
- `When` some action is carried out
- `Then` a particular set of observable consequences should obtain
- `And` is special keyword that can be used to replace Given, When and Then when the statement more than 1 line

### API Feature Grammar

**HTTP Request Testing**

Note that below request is only an example request. If you are confused with the following get request with `POST`, in ordinary `HTTP` request, it should not be.

But if you are interested, this kind behavior is expected in `gRPC`, as `gRPC` is only using `POST` request for all type of request.

```gherkin
Feature: get book detail
    In order eto get book
    As an API user
    I need to be able to request book detail

    Scenario: get book detail
        Given set http-header "Content-Type: application/json"
        And set request body:
            """
            {
                "book_id": 10
            }
            """
        When I send "POST" request to "http://localhost:9000/v1/book/detail"
        Then the response code should be "200"
        And the response header "Content-Type" should be "text/plain; charset=utf-8"
        And the response should match json:
            """
            {
                "book_id": 10,
                "name": "testing",
                "author": "what?"
            }
            """
```

The grammar explanation is as following:

1. `Given set http-header "Content-Type: application/json"` is the sentence to set the http request header
1. `And set request body:` is the sentence to set the http request body
1. `When I send "POST" request to "localhost:9000/v1/book/detail` is the sentence to set the request method and endpoint of a request. At this point of time, the request will be made to the given endpoint.
1. `Then the response code should be "200` is the sentence to check the http code for response.
1. `And the response header "Content-Type" should be "application/json"` is the sentence to check whether certain http header in response is responding with the correct value.
1. `And the response should match json:` is the sentece to check whether the response body is match with what we expect

**Endpoints Mapping**

It is also possible to use endpoints mapping when requesting `api`. But what endpoint mapping means? It means that you should pre-define what endpoint name is available in your test suite.

It will allow the `scenario` to have this kind of grammar:

```gherkin
Feature: get book detail
    In order eto get book
    As an API user
    I need to be able to request book detail

    Scenario: get book detail with endpoint mapping and path
        Given set http-header "Content-Type: application/json"
        And set request body:
            """
            {
                "book_id": 10
            }
            """
        When I send "POST" request to "book service" with path "/v1/book/detail"
        Then the response code should be "200"
        And the response header "Content-Type" should be "text/plain; charset=utf-8"
        And the response should match json:
            """
            {
                "book_id": 10,
                "name": "testing",
                "author": "what?"
            }
            """
```

Look at  `When I send "POST" request to "book service" with path "/v1/book/detail"` sentence, it is different from the previous example. With endpoints mapping, you will be able to use `"{service_name} service"` grammar instead stating the raw endpoint.

To enable this feature, you need to pass `APIFeatureOptions` when initializing `APIFeature`. For example:

```go
apiFeature := &cucumber.APIFeature{
    Options: cucumber.APIFeatureOptions{
        EndpointMapping: map[string]string{
            "book": "http://127.0.0.1:9863",
        },
    },
}
```

**Scenario Outline For API**

For ease of testing, it is also possible to use `scenario outline` for api. This is the example of it:

```gherkin
Scenario Outline: this is an example of scenario outlines
    Given set request header "<request_content_type>"
    And set request body:
        """
        <request_body>
        """
    When I send "<method>" request to "<service> service" with path "<path>"
    Then the response code should be <response_code>
    And the response header "<response_header_key>" should be "<response_header_value>"
    And the response should match json:
        """
        <response_body>
        """
    Examples:
        | request_content_type           | request_body     | method | service | path            | response_code | response_header_key | response_header_value     | response_body                                         |
        | Content-Type: application/json | {"book_id": 10}  | POST   | book    | /v1/book/detail | 200           | Content-Type        | text/plain; charset=utf-8 | {"book_id": 10, "name": "testing", "author": "what?"} |
        | Content-Type: application/json | {"book_id": 20}  | POST   | book    | /v1/book/detail | 200           | Content-Type        | text/plain; charset=utf-8 | {"book_id": 20, "name": "testing", "author": "what?"} |
```

The table in `examples` will automatically looped when the test is executed.

You can see the example in `cucumber_test.go`

## How To Use The Library

First, you need to create a `folder` named `features` inside your project/working directory.

Then, inside the `features` folder, create a `feature` file, for example `api.feature`:

```
|- features
    |- api.feature
```

After that, you need to write a feature, scenario and steps inside the `api.feature`, just like the example above.

There are two ways of running `features` using `godog`:

1. Create a program with package `main` and invoke `godoc` command 
1. Create a test and invoke the test within `TestMain(t *testing.M)`

### Create A Program With Package Main

Create a go program:

```
|- features
    |- api.feature
|- main.go
```

Inside the `main.go`:

```go
package main

import (
    "github.com/albertwidi/go-project-example/internal/pkg/cucumber"
    "github.com/cucumber/godog"
)

func main() {
}
```

Then create FeatureContext function below our main function:

```go
package main

import (
    "github.com/albertwidi/go-project-example/internal/pkg/cucumber"
    "github.com/cucumber/godog"
)

func main() {
}

// godog command will read this function only
func FeatureContext(s *godog.Suite) {
    c := cucumber.New(nil)
    if err := c.RegisterFeatures(&cucumber.APIFeature{}); err != nil {
        log.Fatal(err)
    }
    c.FeatureContext(s)
}
```

Then invoke `godog --format=pretty` inside your working directory

Note that our `func main()` doesn't do anything. It is because `godoc` command doesn't care about `func main()` and only invoke `FeatureContext(s *godog.Suite)` function.

### Create A Test Using TestMain

Create a test program:

```
|- features
    |- api.feature
|- some_test.go
```

Inside the `some_test.go`:

```go
package some_test

import (
	"github.com/cucumber/godog"
)

func TestMain(m *testing.M) {
	var options = godog.Options{
		Output: os.Stdout,
		Format: "pretty",
	}

	flag.Parse()
	options.Paths = flag.Args()

	c, err := cucumber.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := c.RegisterFeatures(&cucumber.APIFeature{}); err != nil {
        log.Fatal(err)
    }

	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		c.FeatureContext(s)
	}, options)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
```

Then invoke `go test *.go` command inside your working directory

## TODO

- Set grammar restrictions per `step`