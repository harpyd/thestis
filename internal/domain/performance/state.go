package performance

type State string

const (
	NotPerformed State = ""
	Performing   State = "performing"
	Passed       State = "passed"
	Failed       State = "failed"
	Crashed      State = "crashed"
	Canceled     State = "canceled"
)

func (s State) String() string {
	return string(s)
}

type stateTransitionRules map[State]map[State]State

func (r stateTransitionRules) apply(from, to State) State {
	return r[from][to]
}

func newCommonStateTransitionRules() stateTransitionRules {
	return stateTransitionRules{
		NotPerformed: fromNotPerformedTransitionRules(true),
		Performing:   fromPerformingTransitionRules(true),
		Failed:       fromFailedTransitionRules(),
		Crashed:      fromErrorTransitionRules(),
		Canceled:     fromCanceledTransitionRules(),
	}
}

func newSpecificStateTransitionRules() stateTransitionRules {
	return stateTransitionRules{
		NotPerformed: fromNotPerformedTransitionRules(false),
		Performing:   fromPerformingTransitionRules(false),
		Passed:       fromPassedTransitionRules(),
		Failed:       fromFailedTransitionRules(),
		Crashed:      fromErrorTransitionRules(),
		Canceled:     fromCanceledTransitionRules(),
	}
}

func fromNotPerformedTransitionRules(commonState bool) map[State]State {
	rules := map[State]State{
		NotPerformed: NotPerformed,
		Performing:   Performing,
		Passed:       Passed,
		Failed:       Failed,
		Crashed:      Crashed,
		Canceled:     Canceled,
	}

	if commonState {
		rules[Passed] = Performing
	}

	return rules
}

func fromPerformingTransitionRules(commonState bool) map[State]State {
	rules := map[State]State{
		NotPerformed: Performing,
		Performing:   Performing,
		Passed:       Passed,
		Failed:       Failed,
		Crashed:      Crashed,
		Canceled:     Canceled,
	}

	if commonState {
		rules[Passed] = Performing
	}

	return rules
}

func fromPassedTransitionRules() map[State]State {
	return map[State]State{
		NotPerformed: NotPerformed,
		Performing:   Performing,
		Passed:       Passed,
		Failed:       Failed,
		Crashed:      Crashed,
		Canceled:     Canceled,
	}
}

func fromFailedTransitionRules() map[State]State {
	return map[State]State{
		NotPerformed: Failed,
		Performing:   Failed,
		Passed:       Failed,
		Failed:       Failed,
		Crashed:      Crashed,
		Canceled:     Failed,
	}
}

func fromErrorTransitionRules() map[State]State {
	return map[State]State{
		NotPerformed: Crashed,
		Performing:   Crashed,
		Passed:       Crashed,
		Failed:       Crashed,
		Crashed:      Crashed,
		Canceled:     Crashed,
	}
}

func fromCanceledTransitionRules() map[State]State {
	return map[State]State{
		NotPerformed: Canceled,
		Performing:   Canceled,
		Passed:       Canceled,
		Failed:       Canceled,
		Crashed:      Canceled,
		Canceled:     Canceled,
	}
}
