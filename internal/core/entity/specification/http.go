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

var ErrNoHTTPRequest = errors.New("no request")

func (h HTTP) validate() error {
	var w BuildErrorWrapper

	if h.request.IsZero() {
		w.WithError(ErrNoHTTPRequest)
	} else {
		w.WithError(h.request.validate())
	}

	if !h.response.IsZero() {
		w.WithError(h.response.validate())
	}

	return w.Wrap("HTTP")
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
	return copyBody(r.body)
}

func copyBody(body map[string]interface{}) map[string]interface{} {
	if len(body) == 0 {
		return nil
	}

	return deepcopy.StringInterfaceMap(body)
}

func (r HTTPRequest) IsZero() bool {
	return r.method == NoHTTPMethod && r.url == "" &&
		r.contentType == NoContentType && len(r.body) == 0
}

func (r HTTPRequest) validate() error {
	var w BuildErrorWrapper

	if !r.method.IsValid() {
		w.WithError(NewNotAllowedHTTPMethodError(r.method))
	}

	if !r.contentType.IsValid() {
		w.WithError(NewNotAllowedContentTypeError(r.contentType))
	}

	return w.Wrap("request")
}

func (r HTTPResponse) AllowedCodes() []int {
	return copyAllowedCodes(r.allowedCodes)
}

func copyAllowedCodes(codes []int) []int {
	if len(codes) == 0 {
		return nil
	}

	result := make([]int, len(codes))
	copy(result, codes)

	return result
}

func (r HTTPResponse) AllowedContentType() ContentType {
	return r.allowedContentType
}

func (r HTTPResponse) IsZero() bool {
	return r.allowedContentType == NoContentType && len(r.allowedCodes) == 0
}

func (r HTTPResponse) validate() error {
	var w BuildErrorWrapper

	if !r.allowedContentType.IsValid() {
		w.WithError(NewNotAllowedContentTypeError(r.allowedContentType))
	}

	return w.Wrap("response")
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

func (b *HTTPBuilder) Build() HTTP {
	return HTTP{
		request:  b.requestBuilder.Build(),
		response: b.responseBuilder.Build(),
	}
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

func (b *HTTPRequestBuilder) Build() HTTPRequest {
	return HTTPRequest{
		method:      b.method,
		url:         b.url,
		contentType: b.contentType,
		body:        bodyOrNil(b.body),
	}
}

func bodyOrNil(body map[string]interface{}) map[string]interface{} {
	if len(body) == 0 {
		return nil
	}

	return body
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

func (b *HTTPResponseBuilder) Build() HTTPResponse {
	return HTTPResponse{
		allowedCodes:       allowedCodesOrNil(b.allowedCodes),
		allowedContentType: b.allowedContentType,
	}
}

func allowedCodesOrNil(codes []int) []int {
	if len(codes) == 0 {
		return nil
	}

	return codes
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
