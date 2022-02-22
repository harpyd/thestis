package natsio_test

import (
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/adapter/pubsub/natsio"
)

func TestPerformanceCancelSignalBus(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Integration tests are skipped")
	}

	natsConn, err := nats.Connect(nats.DefaultURL)
	require.NoError(t, err)

	natsBus := natsio.NewPerformanceCancelSignalBus(natsConn)

	canceled, err := natsBus.SubscribePerformanceCancel("d54e2b7c-0edb-4367-b819-d166ca0edd9e")
	require.NoError(t, err)
	require.NotNil(t, canceled)

	go func() {
		err = natsBus.PublishPerformanceCancel("d54e2b7c-0edb-4367-b819-d166ca0edd9e")
		require.NoError(t, err)
	}()

	v, ok := <-canceled
	require.Empty(t, v)
	require.False(t, ok)
}
