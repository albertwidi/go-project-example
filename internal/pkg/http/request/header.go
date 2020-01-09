package request

import "net/http"

// hTTPHeader for http request
type hTTPHeader struct {
	headers    []string
	HTTPheader http.Header
}

// Header function
func Header(kv ...string) *hTTPHeader {
	return &hTTPHeader{}
}

// Headers return http.Header
func (h *hTTPHeader) Headers() http.Header {
	// process passed header
	for idx := range h.headers {
		if idx > 0 {
			idx++
			if idx == len(h.headers)-1 {
				break
			}
		}
		h.HTTPheader.Add(h.headers[idx], h.headers[idx+1])
	}
	return h.HTTPheader
}

type hTTPContentType struct {
	header *hTTPHeader
	key    string
}

// ContentType for requesting http header content-type
func (h *hTTPHeader) ContentType() *hTTPContentType {
	return &hTTPContentType{key: "Content-Type"}
}

// ApplicationFormWWWURLEncoded return x-www-form-urlencoded for http header
func (ct *hTTPContentType) ApplicationFormWWWURLEncoded() *hTTPHeader {
	ct.header.HTTPheader.Add(ct.key, "application/x-www-form-urlencoded")
	return ct.header
}

// ApplicationJSON return content type application json for http header
func (ct *hTTPContentType) ApplicationJSON() *hTTPHeader {
	ct.header.HTTPheader.Add(ct.key, "application/json")
	return ct.header
}
