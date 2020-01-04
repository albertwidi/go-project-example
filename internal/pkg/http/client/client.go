package client

import "net/http"

// Wrapper for http client
type Wrapper struct {
	c *http.Client
}

// Wrap htpp client
func Wrap(client *http.Client) *Wrapper {
	w := Wrapper{client}
	return &w
}

// Options for http client
type Options struct {
}

// New http client
func New(options Options) *Wrapper {
	w := Wrapper{
		c: &http.Client{},
	}
	return &w
}

// Get wrap the http client get request
func (w *Wrapper) Get() {

}
