package mock

import "github.com/harpyd/thestis/internal/domain/performance"

type Step struct {
	from  string
	to    string
	state performance.State
}

func NewTransitionStep(state performance.State, from, to string) Step {
	return Step{
		state: state,
		from:  from,
		to:    to,
	}
}

func NewCanceledStep() Step {
	return Step{
		state: performance.Canceled,
	}
}

func (s Step) FromTo() (from, to string, ok bool) {
	if s.from == "" && s.to == "" {
		return "", "", true
	}

	return s.from, s.to, true
}

func (s Step) State() performance.State {
	return s.state
}

func (s Step) Err() error {
	return nil
}

func (s Step) String() string {
	return "Mocked step"
}
