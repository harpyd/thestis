package specification_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func buildHTTP(
	t *testing.T,
	prepare func(b *specification.HTTPBuilder),
) specification.HTTP {
	t.Helper()

	builder := specification.NewHTTPBuilder()

	prepare(builder)

	http, err := builder.Build()
	require.NoError(t, err)

	return http
}

func TestBuildHTTPWithRequest(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare         func(b *specification.HTTPBuilder)
		ExpectedRequest specification.HTTPRequest
	}{
		{
			Prepare:         func(b *specification.HTTPBuilder) {},
			ExpectedRequest: specification.HTTPRequest{},
		},
		{
			Prepare: func(b *specification.HTTPBuilder) {
				b.WithRequest(func(b *specification.HTTPRequestBuilder) {
					b.WithMethod("GET")
				})
			},
			ExpectedRequest: specification.NewHTTPRequestBuilder().
				WithMethod("GET").
				ErrlessBuild(),
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ExpectedRequest, buildHTTP(t, c.Prepare).Request())
		})
	}
}

func TestBuildHTTPWithResponse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare          func(b *specification.HTTPBuilder)
		ExpectedResponse specification.HTTPResponse
	}{
		{
			Prepare:          func(b *specification.HTTPBuilder) {},
			ExpectedResponse: specification.HTTPResponse{},
		},
		{
			Prepare: func(b *specification.HTTPBuilder) {
				b.WithResponse(func(b *specification.HTTPResponseBuilder) {
					b.WithAllowedCodes([]int{200})
					b.WithAllowedContentType("application/json")
				})
			},
			ExpectedResponse: specification.NewHTTPResponseBuilder().
				WithAllowedCodes([]int{200}).
				WithAllowedContentType("application/json").
				ErrlessBuild(),
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ExpectedResponse, buildHTTP(t, c.Prepare).Response())
		})
	}
}

func buildHTTPRequest(
	t *testing.T,
	prepare func(b *specification.HTTPRequestBuilder),
) specification.HTTPRequest {
	t.Helper()

	builder := specification.NewHTTPRequestBuilder()

	prepare(builder)

	req, err := builder.Build()
	require.NoError(t, err)

	return req
}

func TestBuildHTTPRequestWithURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare     func(b *specification.HTTPRequestBuilder)
		ExpectedURL string
	}{
		{
			Prepare:     func(b *specification.HTTPRequestBuilder) {},
			ExpectedURL: "",
		},
		{
			Prepare: func(b *specification.HTTPRequestBuilder) {
				b.WithURL("")
			},
			ExpectedURL: "",
		},
		{
			Prepare: func(b *specification.HTTPRequestBuilder) {
				b.WithURL("https://api.warehouse/v1/hooves")
			},
			ExpectedURL: "https://api.warehouse/v1/hooves",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ExpectedURL, buildHTTPRequest(t, c.Prepare).URL())
		})
	}
}

func TestBuildHTTPRequestWithMethod(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Method      specification.HTTPMethod
		ShouldBeErr bool
	}{
		{
			Name:        "allowed_empty",
			Method:      specification.NoHTTPMethod,
			ShouldBeErr: false,
		},
		{
			Name:        "allowed_GET",
			Method:      specification.GET,
			ShouldBeErr: false,
		},
		{
			Name:        "not_allowed_get",
			Method:      "get",
			ShouldBeErr: true,
		},
		{
			Name:        "allowed_POST",
			Method:      specification.POST,
			ShouldBeErr: false,
		},
		{
			Name:        "not_allowed_pOST",
			Method:      "pOST",
			ShouldBeErr: true,
		},
		{
			Name:        "allowed_PUT",
			Method:      specification.PUT,
			ShouldBeErr: false,
		},
		{
			Name:        "not_allowed_pUt",
			Method:      "pUt",
			ShouldBeErr: true,
		},
		{
			Name:        "allowed_PATCH",
			Method:      specification.PATCH,
			ShouldBeErr: false,
		},
		{
			Name:        "not_allowed_pAtCH",
			Method:      "pAtCH",
			ShouldBeErr: true,
		},
		{
			Name:        "allowed_DELETE",
			Method:      specification.DELETE,
			ShouldBeErr: false,
		},
		{
			Name:        "not_allowed_delete",
			Method:      "delete",
			ShouldBeErr: true,
		},
		{
			Name:        "allowed_OPTIONS",
			Method:      specification.OPTIONS,
			ShouldBeErr: false,
		},
		{
			Name:        "not_allowed_OPtions",
			Method:      "OPtions",
			ShouldBeErr: true,
		},
		{
			Name:        "allowed_TRACE",
			Method:      specification.TRACE,
			ShouldBeErr: false,
		},
		{
			Name:        "not_allowed_tRACE",
			Method:      "tRACE",
			ShouldBeErr: true,
		},
		{
			Name:        "allowed_CONNECT",
			Method:      specification.CONNECT,
			ShouldBeErr: false,
		},
		{
			Name:        "not_allowed_connECT",
			Method:      "connECT",
			ShouldBeErr: true,
		},
		{
			Name:        "allowed_HEAD",
			Method:      specification.HEAD,
			ShouldBeErr: false,
		},
		{
			Name:        "not_allowed_head",
			Method:      "head",
			ShouldBeErr: true,
		},
		{
			Name:        "not_allowed_PAST",
			Method:      "PAST",
			ShouldBeErr: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewHTTPRequestBuilder()
			builder.WithMethod(c.Method)

			request, err := builder.Build()

			if c.ShouldBeErr {
				require.True(t, specification.IsNotAllowedHTTPMethodError(err))

				return
			}

			require.NoError(t, err)

			require.Equal(t, c.Method, request.Method())
		})
	}
}

func TestBuildHTTPRequestWithContentType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		ContentType specification.ContentType
		ShouldBeErr bool
	}{
		{
			Name:        "allowed_empty",
			ContentType: specification.NoContentType,
			ShouldBeErr: false,
		},
		{
			Name:        "allowed_application/json",
			ContentType: specification.ApplicationJSON,
			ShouldBeErr: false,
		},
		{
			Name:        "allowed_application/xml",
			ContentType: specification.ApplicationXML,
			ShouldBeErr: false,
		},
		{
			Name:        "not_allowed_content/type",
			ContentType: "content/type",
			ShouldBeErr: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewHTTPRequestBuilder()
			builder.WithContentType(c.ContentType)

			request, err := builder.Build()

			if c.ShouldBeErr {
				require.True(t, specification.IsNotAllowedContentTypeError(err))

				return
			}

			require.NoError(t, err)

			require.Equal(t, c.ContentType, request.ContentType())
		})
	}
}

func TestBuildHTTPRequestWithBody(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare      func(b *specification.HTTPRequestBuilder)
		ExpectedBody map[string]interface{}
	}{
		{
			Prepare:      func(b *specification.HTTPRequestBuilder) {},
			ExpectedBody: nil,
		},
		{
			Prepare: func(b *specification.HTTPRequestBuilder) {
				b.WithBody(map[string]interface{}{})
			},
			ExpectedBody: nil,
		},
		{
			Prepare: func(b *specification.HTTPRequestBuilder) {
				b.WithBody(map[string]interface{}{
					"foo": true,
					"bar": false,
				})
			},
			ExpectedBody: map[string]interface{}{
				"foo": true,
				"bar": false,
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ExpectedBody, buildHTTPRequest(t, c.Prepare).Body())
		})
	}
}

func TestHTTPRequestBodyIsImmutable(t *testing.T) {
	t.Parallel()

	givenBody := map[string]interface{}{
		"foo": 1,
		"bar": 2,
		"baz": 3,
	}

	actualBody := buildHTTPRequest(t, func(b *specification.HTTPRequestBuilder) {
		b.WithBody(givenBody)
	}).Body()

	givenBody["foo"] = 100

	require.Equal(t, actualBody, map[string]interface{}{
		"foo": 1,
		"bar": 2,
		"baz": 3,
	})
}

func buildHTTPResponse(
	t *testing.T,
	prepare func(b *specification.HTTPResponseBuilder),
) specification.HTTPResponse {
	t.Helper()

	builder := specification.NewHTTPResponseBuilder()

	prepare(builder)

	resp, err := builder.Build()
	require.NoError(t, err)

	return resp
}

func TestBuildHTTPResponseWithAllowedCodes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare              func(b *specification.HTTPResponseBuilder)
		ExpectedAllowedCodes []int
	}{
		{
			Prepare:              func(b *specification.HTTPResponseBuilder) {},
			ExpectedAllowedCodes: nil,
		},
		{
			Prepare: func(b *specification.HTTPResponseBuilder) {
				b.WithAllowedCodes([]int{})
			},
			ExpectedAllowedCodes: nil,
		},
		{
			Prepare: func(b *specification.HTTPResponseBuilder) {
				b.WithAllowedCodes([]int{200, 204})
			},
			ExpectedAllowedCodes: []int{200, 204},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.ElementsMatch(t, c.ExpectedAllowedCodes, buildHTTPResponse(t, c.Prepare).AllowedCodes())
		})
	}
}

func TestBuildHTTPResponseWithAllowedContentType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		ContentType specification.ContentType
		ShouldBeErr bool
	}{
		{
			Name:        "allowed_empty",
			ContentType: specification.NoContentType,
			ShouldBeErr: false,
		},
		{
			Name:        "allowed_application/json",
			ContentType: specification.ApplicationJSON,
			ShouldBeErr: false,
		},
		{
			Name:        "allowed_application/xml",
			ContentType: specification.ApplicationXML,
			ShouldBeErr: false,
		},
		{
			Name:        "not_allowed_some/content",
			ContentType: "some/content",
			ShouldBeErr: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewHTTPResponseBuilder()
			builder.WithAllowedContentType(c.ContentType)

			request, err := builder.Build()

			if c.ShouldBeErr {
				require.True(t, specification.IsNotAllowedContentTypeError(err))

				return
			}

			require.NoError(t, err)

			require.Equal(t, c.ContentType, request.AllowedContentType())
		})
	}
}

func TestHTTPMethodIsValid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name          string
		Method        specification.HTTPMethod
		ShouldBeValid bool
	}{
		{
			Name:          "allowed_empty",
			Method:        specification.NoHTTPMethod,
			ShouldBeValid: true,
		},
		{
			Name:          "allowed_GET",
			Method:        specification.GET,
			ShouldBeValid: true,
		},
		{
			Name:          "not_allowed_get",
			Method:        "get",
			ShouldBeValid: false,
		},
		{
			Name:          "allowed_POST",
			Method:        specification.POST,
			ShouldBeValid: true,
		},
		{
			Name:          "not_allowed_post",
			Method:        "post",
			ShouldBeValid: false,
		},
		{
			Name:          "allowed_PUT",
			Method:        specification.PUT,
			ShouldBeValid: true,
		},
		{
			Name:          "not_allowed_put",
			Method:        "put",
			ShouldBeValid: false,
		},
		{
			Name:          "allowed_PATCH",
			Method:        specification.PATCH,
			ShouldBeValid: true,
		},
		{
			Name:          "not_allowed_patch",
			Method:        "patch",
			ShouldBeValid: false,
		},
		{
			Name:          "allowed_DELETE",
			Method:        specification.DELETE,
			ShouldBeValid: true,
		},
		{
			Name:          "not_allowed_delete",
			Method:        "delete",
			ShouldBeValid: false,
		},
		{
			Name:          "allowed_OPTIONS",
			Method:        specification.OPTIONS,
			ShouldBeValid: true,
		},
		{
			Name:          "not_allowed_options",
			Method:        "options",
			ShouldBeValid: false,
		},
		{
			Name:          "allowed_TRACE",
			Method:        specification.TRACE,
			ShouldBeValid: true,
		},
		{
			Name:          "not_allowed_trace",
			Method:        "trace",
			ShouldBeValid: false,
		},
		{
			Name:          "allowed_CONNECT",
			Method:        specification.CONNECT,
			ShouldBeValid: true,
		},
		{
			Name:          "not_allowed_connect",
			Method:        "connect",
			ShouldBeValid: false,
		},
		{
			Name:          "allowed_HEAD",
			Method:        specification.HEAD,
			ShouldBeValid: true,
		},
		{
			Name:          "not_allowed_head",
			Method:        "head",
			ShouldBeValid: false,
		},
		{
			Name:          "not_allowed_PAST",
			Method:        "PAST",
			ShouldBeValid: false,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ShouldBeValid, c.Method.IsValid())
		})
	}
}

func TestContentTypeIsValid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name          string
		ContentType   specification.ContentType
		ShouldBeValid bool
	}{
		{
			Name:          "allowed_empty",
			ContentType:   specification.NoContentType,
			ShouldBeValid: true,
		},
		{
			Name:          "allowed_application/json",
			ContentType:   specification.ApplicationJSON,
			ShouldBeValid: true,
		},
		{
			Name:          "not_allowed_application/JSON",
			ContentType:   "application/JSON",
			ShouldBeValid: false,
		},
		{
			Name:          "allowed_application/xml",
			ContentType:   specification.ApplicationXML,
			ShouldBeValid: true,
		},
		{
			Name:          "not_allowed_application/XML",
			ContentType:   "application/XML",
			ShouldBeValid: false,
		},
		{
			Name:          "not_allowed_some/content",
			ContentType:   "some/content",
			ShouldBeValid: false,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ShouldBeValid, c.ContentType.IsValid())
		})
	}
}

func TestHTTPErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		Err      error
		IsErr    func(err error) bool
		Reversed bool
	}{
		{
			Name:  "build_http_error",
			Err:   specification.NewBuildHTTPError(errors.New("wrong")),
			IsErr: specification.IsBuildHTTPError,
		},
		{
			Name:     "NON_build_http_error",
			Err:      errors.New("wrong"),
			IsErr:    specification.IsBuildHTTPError,
			Reversed: true,
		},
		{
			Name:  "build_http_request_error",
			Err:   specification.NewBuildHTTPRequestError(errors.New("wrong")),
			IsErr: specification.IsBuildHTTPRequestError,
		},
		{
			Name:     "NON_build_http_request_error",
			Err:      errors.New("wrong"),
			IsErr:    specification.IsBuildHTTPRequestError,
			Reversed: true,
		},
		{
			Name:  "build_http_response_error",
			Err:   specification.NewBuildHTTPResponseError(errors.New("something")),
			IsErr: specification.IsBuildHTTPResponseError,
		},
		{
			Name:     "NON_build_http_response_error",
			Err:      errors.New("something"),
			IsErr:    specification.IsBuildHTTPResponseError,
			Reversed: true,
		},
		{
			Name:  "not_allowed_content_type_error",
			Err:   specification.NewNotAllowedContentTypeError("some/content"),
			IsErr: specification.IsNotAllowedContentTypeError,
		},
		{
			Name:     "NON_not_allowed_content_type_error",
			Err:      errors.New("some/content"),
			IsErr:    specification.IsNotAllowedContentTypeError,
			Reversed: true,
		},
		{
			Name:  "not_allowed_http_method_error",
			Err:   specification.NewNotAllowedHTTPMethodError("POZT"),
			IsErr: specification.IsNotAllowedHTTPMethodError,
		},
		{
			Name:     "NON_not_allowed_http_method_error",
			Err:      errors.New("POZT"),
			IsErr:    specification.IsNotAllowedHTTPMethodError,
			Reversed: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			if c.Reversed {
				require.False(t, c.IsErr(c.Err))

				return
			}

			require.True(t, c.IsErr(c.Err))
		})
	}
}
