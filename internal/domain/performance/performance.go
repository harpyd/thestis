package performance

import (
	"sync"

	"github.com/harpyd/thestis/internal/domain/specification"
)

type performerType uint8

const (
	httpPerformer performerType = iota + 1
	assertionPerformer
)

type (
	Performance struct {
		context    *Context
		performers map[performerType]Performer
		graph      [][]action
	}

	Performer interface {
		Perform(c *Context, thesis specification.Thesis)
	}

	action struct {
		thesis        specification.Thesis
		performerType performerType
	}

	Option struct {
		performer     Performer
		performerType performerType
	}
)

func WithHTTPPerformer(performer Performer) Option {
	return Option{
		performer:     performer,
		performerType: httpPerformer,
	}
}

func WithAssertionPerformer(performer Performer) Option {
	return Option{
		performer:     performer,
		performerType: assertionPerformer,
	}
}

func NewFromSpecification(spec *specification.Specification, opts ...Option) *Performance {
	p := &Performance{
		performers: make(map[performerType]Performer, len(opts)),
	}

	// nolint
	// TODO: topological sort of graph

	for _, opt := range opts {
		p.performers[opt.performerType] = opt.performer
	}

	return p
}

func (p *Performance) Start() {
	for i := range p.graph {
		var wg sync.WaitGroup

		wg.Add(len(p.graph[i]))

		for _, a := range p.graph[i] {
			go func(a action) {
				defer wg.Done()

				p.perform(a)
			}(a)
		}

		wg.Wait()
	}
}

func (p *Performance) perform(a action) {
	performer, ok := p.performers[a.performerType]
	if !ok {
		return
	}

	performer.Perform(p.context, a.thesis)
}
