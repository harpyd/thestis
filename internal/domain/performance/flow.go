package performance

type State string

const (
	NotPerformed State = ""
	Performing   State = "performing"
	Passed       State = "passed"
	Failed       State = "failed"
	Error        State = "error"
)

type (
	Flow struct {
		steps map[string]map[string]Step
	}

	Step struct {
		from  string
		to    string
		state State
		fail  error
		err   error
	}
)

func newFlow(graph actionGraph) *Flow {
	steps := make(map[string]map[string]Step, len(graph))

	for from, as := range graph {
		steps[from] = make(map[string]Step, len(as))

		for to := range as {
			steps[from][to] = Step{state: NotPerformed}
		}
	}

	return &Flow{steps: steps}
}

func (f Flow) goToPerforming(from, to string) {
	step := f.steps[from][to]

	step.state = Performing
}

func (f Flow) goToPassed(from, to string) {
	step := f.steps[from][to]

	step.state = Passed
}

func (f Flow) goToFailed(from, to string, fail error) {
	step := f.steps[from][to]

	step.state = Failed
	step.fail = fail
}

func (f Flow) goToError(from, to string, err error) {
	step := f.steps[from][to]

	step.state = Error
	step.err = err
}

func (s Step) From() string {
	return s.from
}

func (s Step) To() string {
	return s.to
}

func (s Step) State() State {
	return s.state
}

func (s Step) Fail() error {
	return s.fail
}

func (s Step) Err() error {
	return s.err
}
