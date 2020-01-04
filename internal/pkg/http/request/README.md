# HTTP Request Builder

This package is an HTTP request builder for go

## Get Request

```go
import "github.com/albertwidi/internal/pkg/http/request"

func main() {
    // req  is a *http.Request 
    req, err := request.New(context.Background()).
        Get("https://google.com").
        Compile()
    if err != nil {
        // do something with error
    }
```

### Get Request With Query

```go
import "github.com/albertwidi/internal/pkg/http/request"

func main() {
    // req  is a *http.Request 
    req, err := request.New(context.Background()).
        Get("https://google.com").
        Query("key", "value").
        Compile()
    if err != nil {
        // do something with error
    }
```

`Query` receive variadic parameters with type `string`

## Post Request

```go
import "github.com/albertwidi/internal/pkg/http/request"

func main() {
    // req  is a *http.Request 
    req, err := request.New(context.Background()).
        Post("https://google.com").
        Compile()
    if err != nil {
        // do something with error
    }
```

### Post Form

```go
import "github.com/albertwidi/internal/pkg/http/request"

func main() {
    // req  is a *http.Request 
    req, err := request.New(context.Background()).
        Post("https://google.com").
        PostForm("key", "value")
        Compile()
    if err != nil {
        // do something with error
    }
```

`PostForm` receive variadic parameters with type `string`

## Request Body

### JSON

```go
import "github.com/albertwidi/internal/pkg/http/request"

func main() {
    s := struct {
        Asd string `json:"asd"`
        Jkl string `json:"jkl"`
    }{}

    // req  is a *http.Request 
    req, err := request.New(context.Background()).
        Post("https://google.com").
        BodyJSON(s).
        Compile()
    if err != nil {
        // do something with error
    }
```

### Raw/io.Reader

```go
func main() {
    s := struct {
        Asd string `json:"asd"`
        Jkl string `json:"jkl"`
    }{}

    // req  is a *http.Request 
    req, err := request.New(context.Background()).
        Post("https://google.com").
        BodyJSON(s).
        Compile()
    if err != nil {
        // do something with error
    }
```

## Version Selection Header

To be added