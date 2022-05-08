package mock

import (
	"context"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/flow"
	"github.com/harpyd/thestis/internal/core/entity/performance"
)

type StepsPolicy struct{}

func NewStepsPolicy() StepsPolicy {
	return StepsPolicy{}
}

func (p StepsPolicy) HandleSteps(
	_ context.Context,
	_ *flow.Flow,
	steps <-chan performance.Step,
	messages chan<- service.Message,
) {
	for s := range steps {
		messages <- service.NewMessageFromStep(s)
	}
}
