package natsio

import (
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/harpyd/thestis/internal/app"
)

type (
	PerformanceCancelSignalBus struct {
		conn *nats.Conn
	}
)

func NewPerformanceCancelSignalBus(conn *nats.Conn) PerformanceCancelSignalBus {
	return PerformanceCancelSignalBus{conn: conn}
}

func (p PerformanceCancelSignalBus) PublishPerformanceCancel(perfID string) error {
	return app.NewPublishCancelError(p.conn.Publish(subject(perfID), []byte{1}))
}

func (p PerformanceCancelSignalBus) SubscribePerformanceCancel(perfID string) (<-chan app.CancelSignal, error) {
	canceled := make(chan app.CancelSignal)

	sub, err := p.conn.Subscribe(subject(perfID), func(msg *nats.Msg) {
		close(canceled)
	})
	if err != nil {
		return nil, app.NewSubscribeCancelError(err)
	}

	if err := sub.AutoUnsubscribe(1); err != nil {
		return nil, app.NewSubscribeCancelError(err)
	}

	return canceled, nil
}

func subject(perfID string) string {
	return fmt.Sprintf("performance.canceled.%s", perfID)
}
