package specification_test

import (
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestHTTPBuilder_WithRequest(t *testing.T) {
	t.Parallel()

	builder := specification.NewHTTPBuilder()
	builder.WithRequest(func(b *specification.HTTPRequestBuilder) {
		b.WithMethod("GET")
		b.WithURL("https://api.warehouse/v1/horns")
	})

	http, err := builder.Build()

	require.NoError(t, err)
	expectedRequest, err := specification.NewHTTPRequestBuilder().
		WithMethod("GET").
		WithURL("https://api.warehouse/v1/horns").
		Build()
	require.NoError(t, err)
	require.Equal(t, expectedRequest, http.Request())
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

func TestIsBuildHTTPError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "build_http_error_is_build_http_error",
			Err:       specification.NewBuildHTTPError(errors.New("wrong")),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_build_http_error",
			Err:       errors.New("wrong"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsBuildHTTPError(c.Err))
		})
	}
}

func TestIsBuildHTTPRequestError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "build_http_request_error_is_build_http_request_error",
			Err:       specification.NewBuildHTTPRequestError(errors.New("wrong")),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_build_http_request_error",
			Err:       errors.New("wrong"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsBuildHTTPRequestError(c.Err))
		})
	}
}

func TestIsBuildHTTPResponseError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "build_http_response_error_is_build_http_response_error",
			Err:       specification.NewBuildHTTPResponseError(errors.New("something")),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_build_http_response_error",
			Err:       errors.New("something"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsBuildHTTPResponseError(c.Err))
		})
	}
}

func TestIsNotAllowedContentTypeError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "not_allowed_content_type_error_is_not_allowed_content_type_error",
			Err:       specification.NewNotAllowedContentTypeError("some/content"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_not_allowed_content_type_error",
			Err:       specification.NewNoSuchStoryError("some/content"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsNotAllowedContentTypeError(c.Err))
		})
	}
}

func TestIsNotAllowedHTTPMethodError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "not_allowed_http_method_error_is_not_allowed_http_method_error",
			Err:       specification.NewNotAllowedHTTPMethodError("POZT"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_not_allowed_http_method_error",
			Err:       specification.NewNoSuchThesisError("POZT"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsNotAllowedHTTPMethodError(c.Err))
		})
	}
}
