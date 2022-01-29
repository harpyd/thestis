package app_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/app"
)

func TestIsDatabaseError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "database_error",
			Err:       app.NewDatabaseError(errors.New("failed to connect")),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       errors.New("failed to connect"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, app.IsDatabaseError(c.Err))
		})
	}
}

func TestIsTestCampaignNotFoundError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "test_campaign_not_found_error",
			Err:       app.NewTestCampaignNotFoundError(errors.New("no documents")),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       errors.New("no documents"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, app.IsTestCampaignNotFoundError(c.Err))
		})
	}
}

func TestIsSpecificationNotFoundError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "specification_not_found_error",
			Err:       app.NewSpecificationNotFoundError(errors.New("no documents")),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       errors.New("no documents"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, app.IsSpecificationNotFoundError(c.Err))
		})
	}
}

func TestIsPerformanceNotFoundError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "performance_not_found_error",
			Err:       app.NewPerformanceNotFoundError(errors.New("no documents")),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       errors.New("no documents"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, app.IsPerformanceNotFoundError(c.Err))
		})
	}
}

func TestIsFlowNotFoundError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "flow_not_found_error",
			Err:       app.NewFlowNotFoundError(errors.New("no documents")),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       errors.New("no documents"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, app.IsFlowNotFoundError(c.Err))
		})
	}
}