package performance_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/performance"
)

func TestPerformance_Start(t *testing.T) {
	t.Parallel()
}

func TestIsCyclicPerformanceGraphError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "cyclic_performance_error",
			Err:       performance.NewCyclicPerformanceGraphError("from", "to"),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       errors.New("from to"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, performance.IsCyclicPerformanceGraphError(c.Err))
		})
	}
}
