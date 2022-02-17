package mock

import (
	"context"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/performance"
)

type PerformanceMaintainer struct {
	withErr bool
}

func NewPerformanceMaintainer(withErr bool) PerformanceMaintainer {
	return PerformanceMaintainer{withErr: withErr}
}

func (m PerformanceMaintainer) MaintainPerformance(_ context.Context, _ *performance.Performance) (<-chan app.Message, error) {
	if m.withErr {
		return nil, performance.NewAlreadyStartedError()
	}

	messages := make(chan app.Message)
	close(messages)

	return messages, nil
}
