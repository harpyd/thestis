package performance_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/performance"
)

func TestNotPerform(t *testing.T) {
	t.Parallel()

	res := performance.NotPerform()

	require.Equal(t, performance.NotPerformed, res.State())
	require.NoError(t, res.Err())
}

func TestPass(t *testing.T) {
	t.Parallel()

	res := performance.Pass()

	require.Equal(t, performance.Passed, res.State())
	require.NoError(t, res.Err())
}

func TestFail(t *testing.T) {
	t.Parallel()

	res := performance.Fail(errors.New("fail"))

	require.Equal(t, performance.Failed, res.State())
	require.True(t, performance.IsFailedError(res.Err()))
}

func TestCrash(t *testing.T) {
	t.Parallel()

	res := performance.Crash(errors.New("crash"))

	require.Equal(t, performance.Crashed, res.State())
	require.True(t, performance.IsCrashedError(res.Err()))
}
