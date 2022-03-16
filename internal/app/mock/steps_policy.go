package mock

import (
	"context"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/flow"
	"github.com/harpyd/thestis/internal/domain/performance"
)

type StepsPolicy struct{}

func NewStepsPolicy() StepsPolicy {
	return StepsPolicy{}
}

func (p StepsPolicy) HandleSteps(
	_ context.Context,
	_ *flow.Reducer,
	steps <-chan performance.Step,
	messages chan<- app.Message,
) {
	for s := range steps {
		messages <- app.NewMessageFromStep(s)
	}
}
