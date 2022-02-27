package specification

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

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
		requestBuilder  *HTTPRequestBuilder
		responseBuilder *HTTPResponseBuilder
	}

	HTTPRequestBuilder struct {
		method      string
		url         string
		contentType string
		body        map[string]interface{}
	}

	HTTPResponseBuilder struct {
		allowedCodes       []int
		allowedContentType string
	}

	ContentType string

	HTTPMethod string
)

const (
	UnknownContentType ContentType = "!"
	EmptyContentType   ContentType = ""
	ApplicationJSON    ContentType = "application/json"
	ApplicationXML     ContentType = "application/xml"
)

const (
	UnknownHTTPMethod HTTPMethod = "!"
	EmptyHTTPMethod   HTTPMethod = ""
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
	return deepcopy.StringInterfaceMap(r.body)
}

func (r HTTPRequest) IsZero() bool {
	return r.method == EmptyHTTPMethod && r.url == "" &&
		r.contentType == EmptyContentType && len(r.body) == 0
}

func (r HTTPResponse) AllowedCodes() []int {
	return deepcopy.IntSlice(r.allowedCodes)
}

func (r HTTPResponse) AllowedContentType() ContentType {
	return r.allowedContentType
}

func (r HTTPResponse) IsZero() bool {
	return r.allowedContentType == EmptyContentType && len(r.allowedCodes) == 0
}

func newContentTypeFromString(contentType string) (ContentType, error) {
	switch strings.ToLower(contentType) {
	case "":
		return EmptyContentType, nil
	case "application/json":
		return ApplicationJSON, nil
	case "application/xml":
		return ApplicationXML, nil
	}

	return UnknownContentType, NewNotAllowedContentTypeError(contentType)
}

func (ct ContentType) String() string {
	return string(ct)
}

func newHTTPMethodFromString(method string) (HTTPMethod, error) {
	methods := map[string]HTTPMethod{
		EmptyHTTPMethod.String(): EmptyHTTPMethod,
		GET.String():             GET,
		POST.String():            POST,
		PUT.String():             PUT,
		PATCH.String():           PATCH,
		DELETE.String():          DELETE,
		OPTIONS.String():         OPTIONS,
		TRACE.String():           TRACE,
		CONNECT.String():         CONNECT,
		HEAD.String():            HEAD,
	}

	if m, ok := methods[strings.ToUpper(method)]; ok {
		return m, nil
	}

	return UnknownHTTPMethod, NewNotAllowedHTTPMethodError(method)
}

func (m HTTPMethod) String() string {
	return string(m)
}

func NewHTTPBuilder() *HTTPBuilder {
	return &HTTPBuilder{
		requestBuilder:  NewHTTPRequestBuilder(),
		responseBuilder: NewHTTPResponseBuilder(),
	}
}

func (b *HTTPBuilder) Build() (HTTP, error) {
	request, requestErr := b.requestBuilder.Build()
	response, responseErr := b.responseBuilder.Build()

	return HTTP{
		request:  request,
		response: response,
	}, NewBuildHTTPError(multierr.Combine(requestErr, responseErr))
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
	buildFn(b.requestBuilder)

	return b
}

func (b *HTTPBuilder) WithResponse(buildFn func(b *HTTPResponseBuilder)) *HTTPBuilder {
	b.responseBuilder.Reset()
	buildFn(b.responseBuilder)

	return b
}

func NewHTTPRequestBuilder() *HTTPRequestBuilder {
	return &HTTPRequestBuilder{}
}

func (b *HTTPRequestBuilder) Build() (HTTPRequest, error) {
	method, methodErr := newHTTPMethodFromString(b.method)
	ctype, ctypeErr := newContentTypeFromString(b.contentType)

	return HTTPRequest{
		method:      method,
		url:         b.url,
		contentType: ctype,
		body:        deepcopy.StringInterfaceMap(b.body),
	}, NewBuildHTTPRequestError(multierr.Combine(methodErr, ctypeErr))
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

func (b *HTTPRequestBuilder) WithMethod(method string) *HTTPRequestBuilder {
	b.method = method

	return b
}

func (b *HTTPRequestBuilder) WithURL(url string) *HTTPRequestBuilder {
	b.url = url

	return b
}

func (b *HTTPRequestBuilder) WithContentType(contentType string) *HTTPRequestBuilder {
	b.contentType = contentType

	return b
}

func (b *HTTPRequestBuilder) WithBody(body map[string]interface{}) *HTTPRequestBuilder {
	b.body = body

	return b
}

func NewHTTPResponseBuilder() *HTTPResponseBuilder {
	return &HTTPResponseBuilder{}
}

func (b *HTTPResponseBuilder) Build() (HTTPResponse, error) {
	allowedContentType, err := newContentTypeFromString(b.allowedContentType)

	return HTTPResponse{
		allowedCodes:       deepcopy.IntSlice(b.allowedCodes),
		allowedContentType: allowedContentType,
	}, NewBuildHTTPResponseError(err)
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

func (b *HTTPResponseBuilder) WithAllowedContentType(allowedContentType string) *HTTPResponseBuilder {
	b.allowedContentType = allowedContentType

	return b
}

type (
	buildHTTPError struct {
		err error
	}

	buildHTTPRequestError struct {
		err error
	}

	buildHTTPResponseError struct {
		err error
	}

	notAllowedContentTypeError struct {
		contentType string
	}

	notAllowedHTTPMethodError struct {
		method string
	}
)

func NewBuildHTTPError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(buildHTTPError{
		err: err,
	})
}

func IsBuildHTTPError(err error) bool {
	var target buildHTTPError

	return errors.As(err, &target)
}

func (e buildHTTPError) Cause() error {
	return e.err
}

func (e buildHTTPError) Unwrap() error {
	return e.err
}

func (e buildHTTPError) NestedErrors() []error {
	return multierr.Errors(e.err)
}

func (e buildHTTPError) CommonError() string {
	return "HTTP"
}

func (e buildHTTPError) Error() string {
	return fmt.Sprintf("HTTP: %s", e.err)
}

func NewBuildHTTPRequestError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(buildHTTPRequestError{
		err: err,
	})
}

func IsBuildHTTPRequestError(err error) bool {
	var target buildHTTPRequestError

	return errors.As(err, &target)
}

func (e buildHTTPRequestError) Cause() error {
	return e.err
}

func (e buildHTTPRequestError) Unwrap() error {
	return e.err
}

func (e buildHTTPRequestError) NestedErrors() []error {
	return multierr.Errors(e.err)
}

func (e buildHTTPRequestError) CommonError() string {
	return "request"
}

func (e buildHTTPRequestError) Error() string {
	return fmt.Sprintf("request: %s", e.err)
}

func NewBuildHTTPResponseError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(buildHTTPResponseError{
		err: err,
	})
}

func IsBuildHTTPResponseError(err error) bool {
	var target buildHTTPResponseError

	return errors.As(err, &target)
}

func (e buildHTTPResponseError) Cause() error {
	return e.err
}

func (e buildHTTPResponseError) Unwrap() error {
	return e.err
}

func (e buildHTTPResponseError) NestedErrors() []error {
	return multierr.Errors(e.err)
}

func (e buildHTTPResponseError) CommonError() string {
	return "response"
}

func (e buildHTTPResponseError) Error() string {
	return fmt.Sprintf("response: %s", e.err)
}

func NewNotAllowedContentTypeError(contentType string) error {
	return errors.WithStack(notAllowedContentTypeError{
		contentType: contentType,
	})
}

func IsNotAllowedContentTypeError(err error) bool {
	var target notAllowedContentTypeError

	return errors.As(err, &target)
}

func (e notAllowedContentTypeError) Error() string {
	return fmt.Sprintf("content type `%s` not allowed", e.contentType)
}

func NewNotAllowedHTTPMethodError(method string) error {
	return errors.WithStack(notAllowedHTTPMethodError{
		method: method,
	})
}

func IsNotAllowedHTTPMethodError(err error) bool {
	var target notAllowedHTTPMethodError

	return errors.As(err, &target)
}

func (e notAllowedHTTPMethodError) Error() string {
	return fmt.Sprintf("HTTP method `%s` not allowed", e.method)
}
