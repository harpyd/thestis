package mock

import (
	"context"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/performance"
)

type PerformanceMaintainer struct {
	withErr bool
}

func NewPerformanceMaintainer(withErr bool) PerformanceMaintainer {
	return PerformanceMaintainer{withErr: withErr}
}

func (m PerformanceMaintainer) MaintainPerformance(
	_ context.Context,
	_ *performance.Performance,
) (<-chan service.Message, error) {
	if m.withErr {
		return nil, performance.ErrAlreadyStarted
	}

	messages := make(chan service.Message)
	close(messages)

	return messages, nil
}