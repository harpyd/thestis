package mock

import (
	"context"

	"github.com/harpyd/thestis/internal/core/app"
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
	messages chan<- app.Message,
) {
	for s := range steps {
		messages <- app.NewMessageFromStep(s)
	}
}
