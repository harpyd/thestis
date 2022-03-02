package specification_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

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

			builder := specification.NewHTTPBuilder()

			c.Prepare(builder)

			http, err := builder.Build()
			require.NoError(t, err)

			require.Equal(t, c.ExpectedRequest, http.Request())
		})
	}
}

func TestHTTPBuilder_WithResponse(t *testing.T) {
	t.Parallel()

	builder := specification.NewHTTPBuilder()
	builder.WithResponse(func(b *specification.HTTPResponseBuilder) {
		b.WithAllowedCodes([]int{200})
		b.WithAllowedContentType("application/json")
	})

	http, err := builder.Build()

	require.NoError(t, err)
	expectedResponse, err := specification.NewHTTPResponseBuilder().
		WithAllowedCodes([]int{200}).
		WithAllowedContentType("application/json").
		Build()
	require.NoError(t, err)
	require.Equal(t, expectedResponse, http.Response())
}

func TestHTTPRequestBuilder_WithURL(t *testing.T) {
	t.Parallel()

	builder := specification.NewHTTPRequestBuilder()
	builder.WithURL("https://api.warehouse/v1/hooves")

	request, err := builder.Build()

	require.NoError(t, err)
	require.Equal(t, "https://api.warehouse/v1/hooves", request.URL())
}

func TestHTTPRequestBuilder_WithMethod(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Method      string
		ShouldBeErr bool
	}{
		{
			Name:        "build_with_allowed_empty_method",
			Method:      "",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_get_method",
			Method:      "GET",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_post_method",
			Method:      "POST",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_put_method",
			Method:      "PUT",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_patch_method",
			Method:      "PATCH",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_delete_method",
			Method:      "DELETE",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_options_method",
			Method:      "OPTIONS",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_trace_method",
			Method:      "TRACE",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_connect_method",
			Method:      "CONNECT",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_head_method",
			Method:      "HEAD",
			ShouldBeErr: false,
		},
		{
			Name:        "dont_build_with_not_allowed_past_method",
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
			require.Equal(t, strings.ToUpper(c.Method), request.Method().String())
		})
	}
}

func TestHTTPRequestBuilder_WithContentType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		ContentType string
		ShouldBeErr bool
	}{
		{
			Name:        "build_with_allowed_empty_content_type",
			ContentType: "",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_content_type_application/json",
			ContentType: "application/json",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_content_type_application/xml",
			ContentType: "application/xml",
			ShouldBeErr: false,
		},
		{
			Name:        "dont_build_with_not_allowed_content_type",
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
			require.Equal(t, strings.ToLower(c.ContentType), request.ContentType().String())
		})
	}
}

func TestHTTPRequestBuilder_WithBody(t *testing.T) {
	t.Parallel()

	builder := specification.NewHTTPRequestBuilder()
	builder.WithBody(map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": map[string]interface{}{
			"d": 1,
			"e": map[string]interface{}{
				"f": 3,
			},
		},
	})

	request, err := builder.Build()

	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": map[string]interface{}{
			"d": 1,
			"e": map[string]interface{}{
				"f": 3,
			},
		},
	}, request.Body())
}

func TestHTTPResponseBuilder_WithAllowedCodes(t *testing.T) {
	t.Parallel()

	builder := specification.NewHTTPResponseBuilder()
	builder.WithAllowedCodes([]int{200, 404})

	request, err := builder.Build()

	require.NoError(t, err)
	require.Equal(t, []int{200, 404}, request.AllowedCodes())
}

func TestHTTPResponseBuilder_WithAllowedContentType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		ContentType string
		ShouldBeErr bool
	}{
		{
			Name:        "build_with_allowed_empty_content_type",
			ContentType: "",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_content_type_application/json",
			ContentType: "application/json",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_content_type_application/xml",
			ContentType: "application/xml",
			ShouldBeErr: false,
		},
		{
			Name:        "dont_build_with_not_allowed_content_type",
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
			require.Equal(t, strings.ToLower(c.ContentType), request.AllowedContentType().String())
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
