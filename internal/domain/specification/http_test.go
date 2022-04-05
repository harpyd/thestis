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

	var b specification.HTTPBuilder

	prepare(&b)

	return b.Build()
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
			ExpectedRequest: (&specification.HTTPRequestBuilder{}).
				WithMethod("GET").
				Build(),
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			actualRequest := buildHTTP(t, c.Prepare).Request()

			require.Equal(t, c.ExpectedRequest, actualRequest)
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
			ExpectedResponse: (&specification.HTTPResponseBuilder{}).
				WithAllowedCodes([]int{200}).
				WithAllowedContentType("application/json").
				Build(),
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			actualResponse := buildHTTP(t, c.Prepare).Response()

			require.Equal(t, c.ExpectedResponse, actualResponse)
		})
	}
}

func buildHTTPRequest(
	t *testing.T,
	prepare func(b *specification.HTTPRequestBuilder),
) specification.HTTPRequest {
	t.Helper()

	var b specification.HTTPRequestBuilder

	prepare(&b)

	return b.Build()
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

			actualURL := buildHTTPRequest(t, c.Prepare).URL()

			require.Equal(t, c.ExpectedURL, actualURL)
		})
	}
}

func TestBuildHTTPRequestWithMethod(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare        func(b *specification.HTTPRequestBuilder)
		ExpectedMethod specification.HTTPMethod
	}{
		{
			Prepare:        func(b *specification.HTTPRequestBuilder) {},
			ExpectedMethod: specification.NoHTTPMethod,
		},
		{
			Prepare: func(b *specification.HTTPRequestBuilder) {
				b.WithMethod(specification.NoHTTPMethod)
			},
			ExpectedMethod: specification.NoHTTPMethod,
		},
		{
			Prepare: func(b *specification.HTTPRequestBuilder) {
				b.WithMethod(specification.GET)
			},
			ExpectedMethod: specification.GET,
		},
		{
			Prepare: func(b *specification.HTTPRequestBuilder) {
				b.WithMethod(specification.POST)
			},
			ExpectedMethod: specification.POST,
		},
		{
			Prepare: func(b *specification.HTTPRequestBuilder) {
				b.WithMethod(specification.PUT)
			},
			ExpectedMethod: specification.PUT,
		},
		{
			Prepare: func(b *specification.HTTPRequestBuilder) {
				b.WithMethod(specification.PATCH)
			},
			ExpectedMethod: specification.PATCH,
		},
		{
			Prepare: func(b *specification.HTTPRequestBuilder) {
				b.WithMethod(specification.DELETE)
			},
			ExpectedMethod: specification.DELETE,
		},
		{
			Prepare: func(b *specification.HTTPRequestBuilder) {
				b.WithMethod(specification.OPTIONS)
			},
			ExpectedMethod: specification.OPTIONS,
		},
		{
			Prepare: func(b *specification.HTTPRequestBuilder) {
				b.WithMethod(specification.TRACE)
			},
			ExpectedMethod: specification.TRACE,
		},
		{
			Prepare: func(b *specification.HTTPRequestBuilder) {
				b.WithMethod(specification.CONNECT)
			},
			ExpectedMethod: specification.CONNECT,
		},
		{
			Prepare: func(b *specification.HTTPRequestBuilder) {
				b.WithMethod(specification.HEAD)
			},
			ExpectedMethod: specification.HEAD,
		},
		{
			Prepare: func(b *specification.HTTPRequestBuilder) {
				b.WithMethod("UNKNOWN")
			},
			ExpectedMethod: "UNKNOWN",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			actualMethod := buildHTTPRequest(t, c.Prepare).Method()

			require.Equal(t, c.ExpectedMethod, actualMethod)
		})
	}
}

func TestBuildHTTPRequestWithContentType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare             func(b *specification.HTTPRequestBuilder)
		ExpectedContentType specification.ContentType
	}{
		{
			Prepare:             func(b *specification.HTTPRequestBuilder) {},
			ExpectedContentType: specification.NoContentType,
		},
		{
			Prepare: func(b *specification.HTTPRequestBuilder) {
				b.WithContentType(specification.NoContentType)
			},
			ExpectedContentType: specification.NoContentType,
		},
		{
			Prepare: func(b *specification.HTTPRequestBuilder) {
				b.WithContentType(specification.ApplicationJSON)
			},
			ExpectedContentType: specification.ApplicationJSON,
		},
		{
			Prepare: func(b *specification.HTTPRequestBuilder) {
				b.WithContentType(specification.ApplicationXML)
			},
			ExpectedContentType: specification.ApplicationXML,
		},
		{
			Prepare: func(b *specification.HTTPRequestBuilder) {
				b.WithContentType("content/type")
			},
			ExpectedContentType: "content/type",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			actualContentType := buildHTTPRequest(t, c.Prepare).ContentType()

			require.Equal(t, c.ExpectedContentType, actualContentType)
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

			actualBody := buildHTTPRequest(t, c.Prepare).Body()

			require.Equal(t, c.ExpectedBody, actualBody)
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

	var b specification.HTTPResponseBuilder

	prepare(&b)

	return b.Build()
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

			actualAllowedCodes := buildHTTPResponse(t, c.Prepare).AllowedCodes()

			require.ElementsMatch(t, c.ExpectedAllowedCodes, actualAllowedCodes)
		})
	}
}

func TestBuildHTTPResponseWithAllowedContentType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare             func(b *specification.HTTPResponseBuilder)
		ExpectedContentType specification.ContentType
	}{
		{
			Prepare:             func(b *specification.HTTPResponseBuilder) {},
			ExpectedContentType: specification.NoContentType,
		},
		{
			Prepare: func(b *specification.HTTPResponseBuilder) {
				b.WithAllowedContentType(specification.NoContentType)
			},
			ExpectedContentType: specification.NoContentType,
		},
		{
			Prepare: func(b *specification.HTTPResponseBuilder) {
				b.WithAllowedContentType(specification.ApplicationJSON)
			},
			ExpectedContentType: specification.ApplicationJSON,
		},
		{
			Prepare: func(b *specification.HTTPResponseBuilder) {
				b.WithAllowedContentType(specification.ApplicationXML)
			},
			ExpectedContentType: specification.ApplicationXML,
		},
		{
			Prepare: func(b *specification.HTTPResponseBuilder) {
				b.WithAllowedContentType("some/type")
			},
			ExpectedContentType: "some/type",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			actualContentType := buildHTTPResponse(t, c.Prepare).AllowedContentType()

			require.Equal(t, c.ExpectedContentType, actualContentType)
		})
	}
}

func TestHTTPMethodIsValid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Method        specification.HTTPMethod
		ShouldBeValid bool
	}{
		{
			Method:        specification.NoHTTPMethod,
			ShouldBeValid: true,
		},
		{
			Method:        specification.GET,
			ShouldBeValid: true,
		},
		{
			Method:        "get",
			ShouldBeValid: false,
		},
		{
			Method:        specification.POST,
			ShouldBeValid: true,
		},
		{
			Method:        "post",
			ShouldBeValid: false,
		},
		{
			Method:        specification.PUT,
			ShouldBeValid: true,
		},
		{
			Method:        "put",
			ShouldBeValid: false,
		},
		{
			Method:        specification.PATCH,
			ShouldBeValid: true,
		},
		{
			Method:        "patch",
			ShouldBeValid: false,
		},
		{
			Method:        specification.DELETE,
			ShouldBeValid: true,
		},
		{
			Method:        "delete",
			ShouldBeValid: false,
		},
		{
			Method:        specification.OPTIONS,
			ShouldBeValid: true,
		},
		{
			Method:        "options",
			ShouldBeValid: false,
		},
		{
			Method:        specification.TRACE,
			ShouldBeValid: true,
		},
		{
			Method:        "trace",
			ShouldBeValid: false,
		},
		{
			Method:        specification.CONNECT,
			ShouldBeValid: true,
		},
		{
			Method:        "connect",
			ShouldBeValid: false,
		},
		{
			Method:        specification.HEAD,
			ShouldBeValid: true,
		},
		{
			Method:        "head",
			ShouldBeValid: false,
		},
		{
			Method:        "PAST",
			ShouldBeValid: false,
		},
		{
			Method:        specification.UnknownHTTPMethod,
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
		ContentType   specification.ContentType
		ShouldBeValid bool
	}{
		{
			ContentType:   specification.NoContentType,
			ShouldBeValid: true,
		},
		{
			ContentType:   specification.ApplicationJSON,
			ShouldBeValid: true,
		},
		{
			ContentType:   "application/JSON",
			ShouldBeValid: false,
		},
		{
			ContentType:   specification.ApplicationXML,
			ShouldBeValid: true,
		},
		{
			ContentType:   "application/XML",
			ShouldBeValid: false,
		},
		{
			ContentType:   "some/content",
			ShouldBeValid: false,
		},
		{
			ContentType:   specification.UnknownContentType,
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

func TestAsNotAllowedContentTypeError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError          error
		ShouldBeWrapped     bool
		ExpectedContentType specification.ContentType
	}{
		{
			GivenError:      nil,
			ShouldBeWrapped: false,
		},
		{
			GivenError: specification.NewNotAllowedContentTypeError(
				specification.ApplicationJSON,
			),
			ShouldBeWrapped:     true,
			ExpectedContentType: specification.ApplicationJSON,
		},
		{
			GivenError: specification.NewNotAllowedContentTypeError(
				specification.UnknownContentType,
			),
			ShouldBeWrapped:     true,
			ExpectedContentType: specification.UnknownContentType,
		},
		{
			GivenError:          specification.NewNotAllowedContentTypeError("foo"),
			ShouldBeWrapped:     true,
			ExpectedContentType: "foo",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var target *specification.NotAllowedContentTypeError

			if !c.ShouldBeWrapped {
				t.Run("not", func(t *testing.T) {
					require.False(t, errors.As(c.GivenError, &target))
				})

				return
			}

			t.Run("as", func(t *testing.T) {
				require.ErrorAs(t, c.GivenError, &target)

				t.Run("content_type", func(t *testing.T) {
					require.Equal(t, c.ExpectedContentType, target.ContentType())
				})
			})
		})
	}
}

func TestFormatNotAllowedContentTypeError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError          error
		ExpectedErrorString string
	}{
		{
			GivenError:          &specification.NotAllowedContentTypeError{},
			ExpectedErrorString: `content type "" not allowed`,
		},
		{
			GivenError: specification.NewNotAllowedContentTypeError(
				specification.ApplicationXML,
			),
			ExpectedErrorString: `content type "application/xml" not allowed`,
		},
		{
			GivenError: specification.NewNotAllowedContentTypeError(
				"bad",
			),
			ExpectedErrorString: `content type "bad" not allowed`,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.EqualError(t, c.GivenError, c.ExpectedErrorString)
		})
	}
}

func TestAsNotAllowedHTTPMethodError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError      error
		ShouldBeWrapped bool
		ExpectedMethod  specification.HTTPMethod
	}{
		{
			GivenError:      nil,
			ShouldBeWrapped: false,
		},
		{
			GivenError:      &specification.NotAllowedHTTPMethodError{},
			ShouldBeWrapped: true,
			ExpectedMethod:  specification.NoHTTPMethod,
		},
		{
			GivenError: specification.NewNotAllowedHTTPMethodError(
				specification.UnknownHTTPMethod,
			),
			ShouldBeWrapped: true,
			ExpectedMethod:  specification.UnknownHTTPMethod,
		},
		{
			GivenError: specification.NewNotAllowedHTTPMethodError(
				"wrong",
			),
			ShouldBeWrapped: true,
			ExpectedMethod:  "wrong",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var target *specification.NotAllowedHTTPMethodError

			if !c.ShouldBeWrapped {
				t.Run("not", func(t *testing.T) {
					require.False(t, errors.As(c.GivenError, &target))
				})

				return
			}

			t.Run("as", func(t *testing.T) {
				require.ErrorAs(t, c.GivenError, &target)

				t.Run("method", func(t *testing.T) {
					require.Equal(t, c.ExpectedMethod, target.Method())
				})
			})
		})
	}
}

func TestFormatNotAllowedHTTPMethodError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError          error
		ExpectedErrorString string
	}{
		{
			GivenError:          &specification.NotAllowedHTTPMethodError{},
			ExpectedErrorString: `HTTP method "" not allowed`,
		},
		{
			GivenError: specification.NewNotAllowedHTTPMethodError(
				specification.GET,
			),
			ExpectedErrorString: `HTTP method "GET" not allowed`,
		},
		{
			GivenError: specification.NewNotAllowedHTTPMethodError(
				"boom",
			),
			ExpectedErrorString: `HTTP method "boom" not allowed`,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.EqualError(t, c.GivenError, c.ExpectedErrorString)
		})
	}
}
