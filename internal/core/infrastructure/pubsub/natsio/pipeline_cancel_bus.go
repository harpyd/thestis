package natsio

import (
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/harpyd/thestis/internal/core/app/service"
)

type (
	PipelineCancelSignalBus struct {
		conn *nats.Conn
	}
)

func NewPipelineCancelSignalBus(conn *nats.Conn) PipelineCancelSignalBus {
	return PipelineCancelSignalBus{conn: conn}
}

func (p PipelineCancelSignalBus) PublishPipelineCancel(pipeID string) error {
	return service.WrapWithPublishCancelError(p.conn.Publish(subject(pipeID), []byte{}))
}

func (p PipelineCancelSignalBus) SubscribePipelineCancel(pipeID string) (<-chan service.CancelSignal, error) {
	canceled := make(chan service.CancelSignal)

	sub, err := p.conn.Subscribe(subject(pipeID), func(msg *nats.Msg) {
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

func subject(pipeID string) string {
	return fmt.Sprintf("pipeline.canceled.%s", pipeID)
}
