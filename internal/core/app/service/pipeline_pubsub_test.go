package service_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/app/service"
)

func TestAsPublishCancelError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError        error
		ShouldBeWrapped   bool
		ExpectedUnwrapped error
	}{
		{
			GivenError:      nil,
			ShouldBeWrapped: false,
		},
		{
			GivenError:      service.WrapWithPublishCancelError(nil),
			ShouldBeWrapped: false,
		},
		{
			GivenError:        &service.PublishCancelError{},
			ShouldBeWrapped:   true,
			ExpectedUnwrapped: nil,
		},
		{
			GivenError:        service.WrapWithPublishCancelError(errors.New("foo")),
			ShouldBeWrapped:   true,
			ExpectedUnwrapped: errors.New("foo"),
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var target *service.PublishCancelError

			if !c.ShouldBeWrapped {
				t.Run("not", func(t *testing.T) {
					require.False(t, errors.As(c.GivenError, &target))
				})

				return
			}

			t.Run("as", func(t *testing.T) {
				require.ErrorAs(t, c.GivenError, &target)

				t.Run("unwrap", func(t *testing.T) {
					if c.ExpectedUnwrapped != nil {
						require.EqualError(t, target.Unwrap(), c.ExpectedUnwrapped.Error())

						return
					}

					require.NoError(t, target.Unwrap())
				})
			})
		})
	}
}

func TestFormatPublishCancelError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError          error
		ExpectedErrorString string
	}{
		{
			GivenError:          &service.PublishCancelError{},
			ExpectedErrorString: "",
		},
		{
			GivenError:          service.WrapWithPublishCancelError(errors.New("failed")),
			ExpectedErrorString: "publish cancel: failed",
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

func TestAsSubscribeCancelError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError        error
		ShouldBeWrapped   bool
		ExpectedUnwrapped error
	}{
		{
			GivenError:      nil,
			ShouldBeWrapped: false,
		},
		{
			GivenError:      service.WrapWithSubscribeCancelError(nil),
			ShouldBeWrapped: false,
		},
		{
			GivenError:        &service.SubscribeCancelError{},
			ShouldBeWrapped:   true,
			ExpectedUnwrapped: nil,
		},
		{
			GivenError:        service.WrapWithSubscribeCancelError(errors.New("qoo")),
			ShouldBeWrapped:   true,
			ExpectedUnwrapped: errors.New("qoo"),
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var target *service.SubscribeCancelError

			if !c.ShouldBeWrapped {
				t.Run("not", func(t *testing.T) {
					require.False(t, errors.As(c.GivenError, &target))
				})

				return
			}

			t.Run("as", func(t *testing.T) {
				require.ErrorAs(t, c.GivenError, &target)

				t.Run("unwrap", func(t *testing.T) {
					if c.ExpectedUnwrapped != nil {
						require.EqualError(t, target.Unwrap(), c.ExpectedUnwrapped.Error())

						return
					}

					require.NoError(t, target.Unwrap())
				})
			})
		})
	}
}

func TestFormatSubscribeCancelError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError          error
		ExpectedErrorString string
	}{
		{
			GivenError:          &service.SubscribeCancelError{},
			ExpectedErrorString: "",
		},
		{
			GivenError:          service.WrapWithSubscribeCancelError(errors.New("wrong")),
			ExpectedErrorString: "subscribe cancel: wrong",
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
