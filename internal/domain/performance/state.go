package performance

type Event string

const (
	NoEvent      Event = ""
	FiredPerform Event = "perform"
	FiredPass    Event = "pass"
	FiredFail    Event = "fail"
	FiredCrash   Event = "crash"
	FiredCancel  Event = "cancel"
)

func (e Event) String() string {
	return string(e)
}

type State string

const (
	NoState      State = ""
	NotPerformed State = "not performed"
	Performing   State = "performing"
	Passed       State = "passed"
	Failed       State = "failed"
	Crashed      State = "crashed"
	Canceled     State = "canceled"
)

type stateTransitionRules map[State]map[Event]State

func rules() stateTransitionRules {
	return stateTransitionRules{
		NotPerformed: {
			FiredPerform: Performing,
			FiredPass:    Passed,
			FiredFail:    Failed,
			FiredCrash:   Crashed,
			FiredCancel:  Canceled,
		},
		Performing: {
			FiredPerform: Performing,
			FiredPass:    Passed,
			FiredFail:    Failed,
			FiredCrash:   Crashed,
			FiredCancel:  Canceled,
		},
		Passed: {
			FiredPerform: Passed,
			FiredPass:    Passed,
			FiredFail:    Failed,
			FiredCrash:   Crashed,
			FiredCancel:  Passed,
		},
		Failed: {
			FiredPerform: Failed,
			FiredPass:    Failed,
			FiredFail:    Failed,
			FiredCrash:   Crashed,
			FiredCancel:  Failed,
		},
		Crashed: {
			FiredPerform: Crashed,
			FiredPass:    Crashed,
			FiredFail:    Crashed,
			FiredCrash:   Crashed,
			FiredCancel:  Crashed,
		},
		Canceled: {
			FiredPerform: Canceled,
			FiredPass:    Canceled,
			FiredFail:    Canceled,
			FiredCrash:   Canceled,
			FiredCancel:  Canceled,
		},
	}
}

func (s State) Next(with Event) State {
	ss, ok := rules()[s]
	if !ok {
		return NoState
	}

	res, ok := ss[with]
	if !ok {
		return NoState
	}

	return res
}

func (s State) String() string {
	return string(s)
}
