package natsio

import (
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/harpyd/thestis/internal/core/app/service"
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
	return service.WrapWithPublishCancelError(p.conn.Publish(subject(perfID), []byte{}))
}

func (p PerformanceCancelSignalBus) SubscribePerformanceCancel(perfID string) (<-chan service.CancelSignal, error) {
	canceled := make(chan service.CancelSignal)

	sub, err := p.conn.Subscribe(subject(perfID), func(msg *nats.Msg) {
		close(canceled)
	})
	if err != nil {
		return nil, service.WrapWithSubscribeCancelError(err)
	}

	if err := sub.AutoUnsubscribe(1); err != nil {
		return nil, service.WrapWithSubscribeCancelError(err)
	}

	return canceled, nil
}

func subject(perfID string) string {
	return fmt.Sprintf("performance.canceled.%s", perfID)
}
