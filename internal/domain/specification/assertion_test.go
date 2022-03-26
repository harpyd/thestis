package specification_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestBuildAssertionWithMethod(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		GivenMethod specification.AssertionMethod
		ShouldBeErr bool
	}{
		{
			Name:        "allowed_empty",
			GivenMethod: specification.NoAssertionMethod,
			ShouldBeErr: false,
		},
		{
			Name:        "allowed_JSONPATH",
			GivenMethod: "JSONPATH",
			ShouldBeErr: true,
		},
		{
			Name:        "allowed_jsonpath",
			GivenMethod: specification.JSONPath,
			ShouldBeErr: false,
		},
		{
			Name:        "allowed_jsonPATH",
			GivenMethod: "jsonPATH",
			ShouldBeErr: true,
		},
		{
			Name:        "not_allowed_JAYZ",
			GivenMethod: "JAYZ",
			ShouldBeErr: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var b specification.AssertionBuilder

			b.WithMethod(c.GivenMethod)

			assertion, err := b.Build()

			if c.ShouldBeErr {
				var target *specification.NotAllowedAssertionMethodError

				require.ErrorAs(t, err, &target)

				return
			}

			require.NoError(t, err)

			require.Equal(t, c.GivenMethod, assertion.Method())
		})
	}
}

func TestBuildAssertionWithAsserts(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare         func(b *specification.AssertionBuilder)
		ExpectedAsserts []specification.Assert
	}{
		{
			Prepare:         func(b *specification.AssertionBuilder) {},
			ExpectedAsserts: nil,
		},
		{
			Prepare: func(b *specification.AssertionBuilder) {
				b.WithAssert("foo.bar.baz", "somebody")
			},
			ExpectedAsserts: []specification.Assert{
				specification.NewAssert("foo.bar.baz", "somebody"),
			},
		},
		{
			Prepare: func(b *specification.AssertionBuilder) {
				b.
					WithAssert("getSomeBody.response.body.type", "product").
					WithAssert("getSomeBody.response.body.items..price", []int{2100, 1100}).
					WithAssert("getSomeBody.response.body.items..amount", []int{10, 33})
			},
			ExpectedAsserts: []specification.Assert{
				specification.NewAssert("getSomeBody.response.body.type", "product"),
				specification.NewAssert("getSomeBody.response.body.items..price", []int{2100, 1100}),
				specification.NewAssert("getSomeBody.response.body.items..amount", []int{10, 33}),
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var b specification.AssertionBuilder

			c.Prepare(&b)

			assertion, err := b.Build()
			require.NoError(t, err)

			asserts := assertion.Asserts()

			require.ElementsMatch(t, c.ExpectedAsserts, asserts)
		})
	}
}

func TestAssertionMethodIsValid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenMethod   specification.AssertionMethod
		ShouldBeValid bool
	}{
		{
			GivenMethod:   specification.NoAssertionMethod,
			ShouldBeValid: true,
		},
		{
			GivenMethod:   specification.UnknownAssertionMethod,
			ShouldBeValid: false,
		},
		{
			GivenMethod:   specification.JSONPath,
			ShouldBeValid: true,
		},
		{
			GivenMethod:   "JSONpath",
			ShouldBeValid: false,
		},
		{
			GivenMethod:   "somethingelse",
			ShouldBeValid: false,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ShouldBeValid, c.GivenMethod.IsValid())
		})
	}
}

func TestAssert(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenAssert specification.Assert
		Actual      string
		Expected    interface{}
	}{
		{
			GivenAssert: specification.NewAssert("", struct{}{}),
			Actual:      "",
			Expected:    struct{}{},
		},
		{
			GivenAssert: specification.NewAssert("some", "foo"),
			Actual:      "some",
			Expected:    "foo",
		},
		{
			GivenAssert: specification.NewAssert("map", map[string]interface{}{
				"foo": true,
				"bar": false,
			}),
			Actual: "map",
			Expected: map[string]interface{}{
				"foo": true,
				"bar": false,
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			t.Run("actual", func(t *testing.T) {
				require.Equal(t, c.Actual, c.GivenAssert.Actual())
			})

			t.Run("expected", func(t *testing.T) {
				require.Equal(t, c.Expected, c.GivenAssert.Expected())
			})
		})
	}
}

func TestAsNotAllowedAssertionMethodError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError      error
		ShouldBeWrapped bool
		ExpectedMethod  specification.AssertionMethod
	}{
		{
			GivenError:      nil,
			ShouldBeWrapped: false,
		},
		{
			GivenError:      &specification.NotAllowedAssertionMethodError{},
			ShouldBeWrapped: true,
			ExpectedMethod:  specification.NoAssertionMethod,
		},
		{
			GivenError: specification.NewNotAllowedAssertionMethodError(
				specification.JSONPath,
			),
			ShouldBeWrapped: true,
			ExpectedMethod:  specification.JSONPath,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var nerr *specification.NotAllowedAssertionMethodError

			if !c.ShouldBeWrapped {
				t.Run("not", func(t *testing.T) {
					require.False(t, errors.As(c.GivenError, &nerr))
				})

				return
			}

			t.Run("as", func(t *testing.T) {
				require.ErrorAs(t, c.GivenError, &nerr)

				t.Run("method", func(t *testing.T) {
					require.Equal(t, c.ExpectedMethod, nerr.Method())
				})
			})
		})
	}
}

func TestFormatNotAllowedAssertionMethodError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError          error
		ExpectedErrorString string
	}{
		{
			GivenError:          &specification.NotAllowedAssertionMethodError{},
			ExpectedErrorString: "assertion method `` not allowed",
		},
		{
			GivenError: specification.NewNotAllowedAssertionMethodError(
				specification.NoAssertionMethod,
			),
			ExpectedErrorString: "assertion method `` not allowed",
		},
		{
			GivenError: specification.NewNotAllowedAssertionMethodError(
				specification.JSONPath,
			),
			ExpectedErrorString: "assertion method `jsonpath` not allowed",
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
