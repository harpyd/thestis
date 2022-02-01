package app

import "github.com/harpyd/thestis/internal/domain/performance"

type (
	PerformerOption struct {
		performerType performance.PerformerType
		performer     performance.Performer
	}

	PerformerOptions []PerformerOption
)

func WithHTTPPerformer(performer performance.Performer) PerformerOption {
	return PerformerOption{
		performerType: performance.HTTPPerformer,
		performer:     performer,
	}
}

func WithAssertionPerformer(performer performance.Performer) PerformerOption {
	return PerformerOption{
		performerType: performance.AssertionPerformer,
		performer:     performer,
	}
}

func (o PerformerOption) ToPerformanceOption() performance.Option {
	perfOpts := map[performance.PerformerType]performance.Option{
		performance.UnknownPerformer:   func(p *performance.Performance) {},
		performance.EmptyPerformer:     func(p *performance.Performance) {},
		performance.HTTPPerformer:      performance.WithHTTP(o.performer),
		performance.AssertionPerformer: performance.WithAssertion(o.performer),
	}

	opt, ok := perfOpts[o.performerType]
	if !ok {
		return func(p *performance.Performance) {}
	}

	return opt
}

func (os PerformerOptions) ToPerformanceOptions() []performance.Option {
	perfOpts := make([]performance.Option, 0, len(os))

	for _, o := range os {
		perfOpts = append(perfOpts, o.ToPerformanceOption())
	}

	return perfOpts
}
