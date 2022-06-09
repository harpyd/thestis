package pipeline_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/entity/pipeline"
)

func TestConcurrentRWOperationsOnEnvironment(t *testing.T) {
	t.Parallel()

	env := pipeline.NewEnvironment(1)
	env.Store("count", 0)

	const (
		writersNumber = 50
		readersNumber = 100
	)

	require.NotPanics(t, func() {
		var wg sync.WaitGroup

		wg.Add(writersNumber)
		for w := 1; w <= writersNumber; w++ {
			go func(w int) {
				defer wg.Done()

				env.Store("count", w)
			}(w)
		}

		wg.Add(readersNumber)
		for r := 1; r <= readersNumber; r++ {
			go func() {
				defer wg.Done()

				_, _ = env.Load("count")
			}()
		}

		wg.Wait()
	})
}
