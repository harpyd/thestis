package specification_test

import (
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestAssertionBuilder_WithMethod(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Method      string
		ShouldBeErr bool
	}{
		{
			Name:        "build_with_allowed_empty_assertion_method",
			Method:      "",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_jsonpath_assertion_method",
			Method:      "JSONPATH",
			ShouldBeErr: false,
		},
		{
			Name:        "dont_build_with_not_allowed_assertion_method",
			Method:      "JAYZ",
			ShouldBeErr: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewAssertionBuilder()
			builder.WithMethod(c.Method)

			assertion, err := builder.Build()

			if c.ShouldBeErr {
				require.True(t, specification.IsNotAllowedAssertionMethodError(err))

				return
			}

			require.NoError(t, err)
			require.Equal(t, strings.ToLower(c.Method), assertion.Method().String())
		})
	}
}

func TestAssertionBuilder_WithAssert(t *testing.T) {
	t.Parallel()

	builder := specification.NewAssertionBuilder()
	builder.WithAssert("getSomeBody.response.body.type", "product")
	builder.WithAssert("getSomeBody.response.body.items..price", []int{2100, 1100})
	builder.WithAssert("getSomeBody.response.body.items..amount", []int{10, 33})

	assertion, err := builder.Build()

	require.NoError(t, err)

	asserts := assertion.Asserts()

	require.Equal(t, []string{
		"getSomeBody.response.body.type",
		"getSomeBody.response.body.items..price",
		"getSomeBody.response.body.items..amount",
	}, mapAssertsToActual(asserts))

	require.Equal(t, []interface{}{
		"product",
		[]int{2100, 1100},
		[]int{10, 33},
	}, mapAssertsToExpected(asserts))
}

func mapAssertsToActual(asserts []specification.Assert) []string {
	expected := make([]string, 0, len(asserts))
	for _, a := range asserts {
		expected = append(expected, a.Actual())
	}

	return expected
}

func mapAssertsToExpected(asserts []specification.Assert) []interface{} {
	actual := make([]interface{}, 0, len(asserts))
	for _, a := range asserts {
		actual = append(actual, a.Expected())
	}

	return actual
}

func TestIsBuildAssertionError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "build_assertion_error",
			Err:       specification.NewBuildAssertionError(errors.New("some")),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       errors.New("some"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsBuildAssertionError(c.Err))
		})
	}
}

func TestIsNotAllowedAssertionMethodError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "not_allowed_assertion_method_error",
			Err:       specification.NewNotAllowedAssertionMethodError("jzonpad"),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       specification.NewNotAllowedKeywordError("jzonpad"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsNotAllowedAssertionMethodError(c.Err))
		})
	}
}
