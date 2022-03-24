package app_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/app"
)

func TestRepositoryErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		Err      error
		IsErr    func(err error) bool
		Reversed bool
	}{
		{
			Name:  "database_error",
			Err:   app.NewDatabaseError(errors.New("failed to connect")),
			IsErr: app.IsDatabaseError,
		},
		{
			Name:     "NON_database_error",
			Err:      errors.New("failed to connect"),
			IsErr:    app.IsDatabaseError,
			Reversed: true,
		},
		{
			Name:  "test_campaign_not_found_error",
			Err:   app.NewTestCampaignNotFoundError(errors.New("no documents")),
			IsErr: app.IsTestCampaignNotFoundError,
		},
		{
			Name:     "NON_test_campaign_not_found_error",
			Err:      errors.New("no documents"),
			IsErr:    app.IsTestCampaignNotFoundError,
			Reversed: true,
		},
		{
			Name:  "specification_not_found_error",
			Err:   app.NewSpecificationNotFoundError(errors.New("no documents")),
			IsErr: app.IsSpecificationNotFoundError,
		},
		{
			Name:     "NON_specification_not_found_error",
			Err:      errors.New("no documents"),
			IsErr:    app.IsSpecificationNotFoundError,
			Reversed: true,
		},
		{
			Name:  "performance_not_found_error",
			Err:   app.NewPerformanceNotFoundError(errors.New("no documents")),
			IsErr: app.IsPerformanceNotFoundError,
		},
		{
			Name:     "NON_performance_not_found_error",
			Err:      errors.New("no documents"),
			IsErr:    app.IsPerformanceNotFoundError,
			Reversed: true,
		},
		{
			Name:  "flow_not_found_error",
			Err:   app.NewFlowNotFoundError(errors.New("no documents")),
			IsErr: app.IsFlowNotFoundError,
		},
		{
			Name:     "NON_flow_not_found_error",
			Err:      errors.New("no documents"),
			IsErr:    app.IsFlowNotFoundError,
			Reversed: true,
		},
		{
			Name:  "already_exists_error",
			Err:   app.NewAlreadyExistsError(errors.New("duplicate key")),
			IsErr: app.IsAlreadyExistsError,
		},
		{
			Name:     "NON_already_exists_error",
			Err:      errors.New("duplicate key"),
			IsErr:    app.IsAlreadyExistsError,
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
