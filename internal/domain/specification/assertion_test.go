package specification_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestBuildAssertionMethod(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		GivenMethod string
		ShouldBeErr bool
	}{
		{
			Name:        "allowed_empty",
			GivenMethod: "",
			ShouldBeErr: false,
		},
		{
			Name:        "allowed_JSONPATH",
			GivenMethod: "JSONPATH",
			ShouldBeErr: false,
		},
		{
			Name:        "allowed_jsonpath",
			GivenMethod: "jsonpath",
			ShouldBeErr: false,
		},
		{
			Name:        "allowed_jsonPATH",
			GivenMethod: "jsonPATH",
			ShouldBeErr: false,
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

			builder := specification.NewAssertionBuilder()
			builder.WithMethod(c.GivenMethod)

			assertion, err := builder.Build()

			if c.ShouldBeErr {
				require.True(t, specification.IsNotAllowedAssertionMethodError(err))

				return
			}

			require.NoError(t, err)

			require.Equal(t, strings.ToLower(c.GivenMethod), assertion.Method().String())
		})
	}
}

func TestBuildAssertionAsserts(t *testing.T) {
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

			builder := specification.NewAssertionBuilder()

			c.Prepare(builder)

			assertion, err := builder.Build()
			require.NoError(t, err)

			asserts := assertion.Asserts()

			require.Equal(t, c.ExpectedAsserts, asserts)
		})
	}
}

func TestAssertionErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		Err      error
		IsErr    func(err error) bool
		Reversed bool
	}{
		{
			Name:  "build_assertion_error",
			Err:   specification.NewBuildAssertionError(errors.New("foo")),
			IsErr: specification.IsBuildAssertionError,
		},
		{
			Name:     "NON_build_assertion_error",
			Err:      errors.New("foo"),
			IsErr:    specification.IsBuildAssertionError,
			Reversed: true,
		},
		{
			Name:  "not_allowed_assertion_method_error",
			Err:   specification.NewNotAllowedAssertionMethodError("SSD"),
			IsErr: specification.IsNotAllowedAssertionMethodError,
		},
		{
			Name:     "NON_not_allowed_assertion_method_error",
			Err:      errors.New("SSD"),
			IsErr:    specification.IsNotAllowedAssertionMethodError,
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
