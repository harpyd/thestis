package specification

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/pkg/deepcopy"
)

type (
	HTTP struct {
		request  HTTPRequest
		response HTTPResponse
	}

	HTTPRequest struct {
		method      HTTPMethod
		url         string
		contentType ContentType
		body        map[string]interface{}
	}

	HTTPResponse struct {
		allowedCodes       []int
		allowedContentType ContentType
	}

	HTTPBuilder struct {
		requestBuilder  HTTPRequestBuilder
		responseBuilder HTTPResponseBuilder
	}

	HTTPRequestBuilder struct {
		method      HTTPMethod
		url         string
		contentType ContentType
		body        map[string]interface{}
	}

	HTTPResponseBuilder struct {
		allowedCodes       []int
		allowedContentType ContentType
	}

	ContentType string

	HTTPMethod string
)

const (
	UnknownContentType ContentType = "!"
	NoContentType      ContentType = ""
	ApplicationJSON    ContentType = "application/json"
	ApplicationXML     ContentType = "application/xml"
)

const (
	UnknownHTTPMethod HTTPMethod = "!"
	NoHTTPMethod      HTTPMethod = ""
	GET               HTTPMethod = "GET"
	POST              HTTPMethod = "POST"
	PUT               HTTPMethod = "PUT"
	PATCH             HTTPMethod = "PATCH"
	DELETE            HTTPMethod = "DELETE"
	OPTIONS           HTTPMethod = "OPTIONS"
	TRACE             HTTPMethod = "TRACE"
	CONNECT           HTTPMethod = "CONNECT"
	HEAD              HTTPMethod = "HEAD"
)

func (h HTTP) Request() HTTPRequest {
	return h.request
}

func (h HTTP) Response() HTTPResponse {
	return h.response
}

func (h HTTP) IsZero() bool {
	return h.Response().IsZero() && h.Request().IsZero()
}

func (r HTTPRequest) Method() HTTPMethod {
	return r.method
}

func (r HTTPRequest) URL() string {
	return r.url
}

func (r HTTPRequest) ContentType() ContentType {
	return r.contentType
}

func (r HTTPRequest) Body() map[string]interface{} {
	if r.body == nil {
		return nil
	}

	return deepcopy.StringInterfaceMap(r.body)
}

func (r HTTPRequest) IsZero() bool {
	return r.method == NoHTTPMethod && r.url == "" &&
		r.contentType == NoContentType && len(r.body) == 0
}

func (r HTTPResponse) AllowedCodes() []int {
	if len(r.allowedCodes) == 0 {
		return nil
	}

	return deepcopy.IntSlice(r.allowedCodes)
}

func (r HTTPResponse) AllowedContentType() ContentType {
	return r.allowedContentType
}

func (r HTTPResponse) IsZero() bool {
	return r.allowedContentType == NoContentType && len(r.allowedCodes) == 0
}

func (ct ContentType) IsValid() bool {
	switch ct {
	case NoContentType:
		return true
	case ApplicationJSON:
		return true
	case ApplicationXML:
		return true
	case UnknownContentType:
		return false
	}

	return false
}

func (ct ContentType) String() string {
	return string(ct)
}

func (m HTTPMethod) IsValid() bool {
	valid := map[HTTPMethod]bool{
		UnknownHTTPMethod: false,
		NoHTTPMethod:      true,
		GET:               true,
		POST:              true,
		PUT:               true,
		PATCH:             true,
		DELETE:            true,
		OPTIONS:           true,
		TRACE:             true,
		CONNECT:           true,
		HEAD:              true,
	}

	return valid[m]
}

func (m HTTPMethod) String() string {
	return string(m)
}

func (b *HTTPBuilder) Build() (HTTP, error) {
	var w BuildErrorWrapper

	request, err := b.requestBuilder.Build()
	w.WithError(err)

	response, err := b.responseBuilder.Build()
	w.WithError(err)

	return HTTP{
		request:  request,
		response: response,
	}, w.Wrap("HTTP")
}

func (b *HTTPBuilder) ErrlessBuild() HTTP {
	h, _ := b.Build()

	return h
}

func (b *HTTPBuilder) Reset() {
	b.requestBuilder.Reset()
	b.responseBuilder.Reset()
}

func (b *HTTPBuilder) WithRequest(buildFn func(b *HTTPRequestBuilder)) *HTTPBuilder {
	b.requestBuilder.Reset()
	buildFn(&b.requestBuilder)

	return b
}

func (b *HTTPBuilder) WithResponse(buildFn func(b *HTTPResponseBuilder)) *HTTPBuilder {
	b.responseBuilder.Reset()
	buildFn(&b.responseBuilder)

	return b
}

func (b *HTTPRequestBuilder) Build() (HTTPRequest, error) {
	var w BuildErrorWrapper

	if !b.method.IsValid() {
		w.WithError(NewNotAllowedHTTPMethodError(b.method))
	}

	if !b.contentType.IsValid() {
		w.WithError(NewNotAllowedContentTypeError(b.contentType))
	}

	return HTTPRequest{
		method:      b.method,
		url:         b.url,
		contentType: b.contentType,
		body:        copyBody(b.body),
	}, w.Wrap("request")
}

func copyBody(body map[string]interface{}) map[string]interface{} {
	if len(body) == 0 {
		return nil
	}

	return deepcopy.StringInterfaceMap(body)
}

func (b *HTTPRequestBuilder) ErrlessBuild() HTTPRequest {
	r, _ := b.Build()

	return r
}

func (b *HTTPRequestBuilder) Reset() {
	b.method = ""
	b.url = ""
	b.contentType = ""
	b.body = nil
}

func (b *HTTPRequestBuilder) WithMethod(method HTTPMethod) *HTTPRequestBuilder {
	b.method = method

	return b
}

func (b *HTTPRequestBuilder) WithURL(url string) *HTTPRequestBuilder {
	b.url = url

	return b
}

func (b *HTTPRequestBuilder) WithContentType(contentType ContentType) *HTTPRequestBuilder {
	b.contentType = contentType

	return b
}

func (b *HTTPRequestBuilder) WithBody(body map[string]interface{}) *HTTPRequestBuilder {
	b.body = body

	return b
}

func (b *HTTPResponseBuilder) Build() (HTTPResponse, error) {
	var w BuildErrorWrapper

	if !b.allowedContentType.IsValid() {
		w.WithError(NewNotAllowedContentTypeError(b.allowedContentType))
	}

	return HTTPResponse{
		allowedCodes:       copyAllowedCodes(b.allowedCodes),
		allowedContentType: b.allowedContentType,
	}, w.Wrap("response")
}

func copyAllowedCodes(codes []int) []int {
	if len(codes) == 0 {
		return nil
	}

	return deepcopy.IntSlice(codes)
}

func (b *HTTPResponseBuilder) ErrlessBuild() HTTPResponse {
	r, _ := b.Build()

	return r
}

func (b *HTTPResponseBuilder) Reset() {
	b.allowedCodes = nil
	b.allowedContentType = ""
}

func (b *HTTPResponseBuilder) WithAllowedCodes(allowedCodes []int) *HTTPResponseBuilder {
	b.allowedCodes = allowedCodes

	return b
}

func (b *HTTPResponseBuilder) WithAllowedContentType(
	allowedContentType ContentType,
) *HTTPResponseBuilder {
	b.allowedContentType = allowedContentType

	return b
}

type NotAllowedContentTypeError struct {
	contentType ContentType
}

func NewNotAllowedContentTypeError(contentType ContentType) error {
	return errors.WithStack(&NotAllowedContentTypeError{
		contentType: contentType,
	})
}

func (e *NotAllowedContentTypeError) ContentType() ContentType {
	return e.contentType
}

func (e *NotAllowedContentTypeError) Error() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("content type %q not allowed", e.contentType)
}

type NotAllowedHTTPMethodError struct {
	method HTTPMethod
}

func NewNotAllowedHTTPMethodError(method HTTPMethod) error {
	return errors.WithStack(&NotAllowedHTTPMethodError{
		method: method,
	})
}

func (e *NotAllowedHTTPMethodError) Method() HTTPMethod {
	return e.method
}

func (e *NotAllowedHTTPMethodError) Error() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("HTTP method %q not allowed", e.method)
}
