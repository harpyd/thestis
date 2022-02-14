package mock

import (
	"context"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/performance"
)

type FlowManager struct {
	withErr bool
}

func NewFlowManager(withErr bool) FlowManager {
	return FlowManager{withErr: withErr}
}

func (m FlowManager) ManageFlow(_ context.Context, _ *performance.Performance) (<-chan app.Message, error) {
	if m.withErr {
		return nil, performance.NewAlreadyStartedError()
	}

	messages := make(chan app.Message)
	defer close(messages)

	return messages, nil
}
