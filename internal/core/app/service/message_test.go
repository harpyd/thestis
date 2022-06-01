package service_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/performance"
	"github.com/harpyd/thestis/internal/core/entity/specification"
)

func TestMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Message        service.Message
		ExpectedEvent  performance.Event
		ExpectedErr    error
		ExpectedString string
	}{
		{
			Message:        service.Message{},
			ExpectedEvent:  performance.NoEvent,
			ExpectedErr:    nil,
			ExpectedString: "",
		},
		{
			Message:        service.NewMessageFromError(errors.New("foo")),
			ExpectedEvent:  performance.NoEvent,
			ExpectedErr:    errors.New("foo"),
			ExpectedString: "foo",
		},
		{
			Message:        service.NewMessageFromError(nil),
			ExpectedEvent:  performance.NoEvent,
			ExpectedErr:    nil,
			ExpectedString: "",
		},
		{
			Message: service.NewMessageFromStep(performance.NewThesisStep(
				specification.NewThesisSlug("foo", "bar", "baz"),
				performance.HTTPPerformer,
				performance.FiredPerform,
			)),
			ExpectedEvent:  performance.FiredPerform,
			ExpectedErr:    nil,
			ExpectedString: "foo.bar.baz: event = perform, type = HTTP",
		},
		{
			Message: service.NewMessageFromStep(performance.NewScenarioStepWithErr(
				errors.New("something wrong"),
				specification.NewScenarioSlug("foo", "bar"),
				performance.FiredCrash,
			)),
			ExpectedEvent:  performance.FiredCrash,
			ExpectedErr:    errors.New("something wrong"),
			ExpectedString: "foo.bar: event = crash, err = something wrong",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			t.Run("event", func(t *testing.T) {
				require.Equal(t, c.ExpectedEvent, c.Message.Event())
			})

			if c.ExpectedErr != nil {
				t.Run("err", func(t *testing.T) {
					require.EqualError(t, c.Message.Err(), c.ExpectedErr.Error())
				})
			} else {
				t.Run("no_err", func(t *testing.T) {
					require.NoError(t, c.Message.Err())
				})
			}

			t.Run("string", func(t *testing.T) {
				require.Equal(t, c.ExpectedString, c.Message.String())
			})
		})
	}
}
