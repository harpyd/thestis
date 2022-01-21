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
		Failed:       fromFailedTransitionRules(true),
		Crashed:      fromCrashedTransitionRules(true),
		Canceled:     fromCanceledTransitionRules(),
	}
}

func newSpecificStateTransitionRules() stateTransitionRules {
	return stateTransitionRules{
		NotPerformed: fromNotPerformedTransitionRules(false),
		Performing:   fromPerformingTransitionRules(false),
		Passed:       fromPassedTransitionRules(),
		Failed:       fromFailedTransitionRules(false),
		Crashed:      fromCrashedTransitionRules(false),
	}
}

func fromNotPerformedTransitionRules(commonState bool) map[State]State {
	rules := map[State]State{
		NotPerformed: NotPerformed,
		Performing:   Performing,
		Passed:       Passed,
		Failed:       Failed,
		Crashed:      Crashed,
	}

	if commonState {
		rules[Passed] = Performing
		rules[Canceled] = Canceled
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
	}

	if commonState {
		rules[Passed] = Performing
		rules[Canceled] = Canceled
	}

	return rules
}

func fromPassedTransitionRules() map[State]State {
	return map[State]State{
		NotPerformed: Passed,
		Performing:   Passed,
		Passed:       Passed,
		Failed:       Failed,
		Crashed:      Crashed,
		Canceled:     Canceled,
	}
}

func fromFailedTransitionRules(commonState bool) map[State]State {
	rules := map[State]State{
		NotPerformed: Failed,
		Performing:   Failed,
		Passed:       Failed,
		Failed:       Failed,
		Crashed:      Crashed,
	}

	if commonState {
		rules[Canceled] = Failed
	}

	return rules
}

func fromCrashedTransitionRules(commonState bool) map[State]State {
	rules := map[State]State{
		NotPerformed: Crashed,
		Performing:   Crashed,
		Passed:       Crashed,
		Failed:       Crashed,
		Crashed:      Crashed,
	}

	if commonState {
		rules[Canceled] = Crashed
	}

	return rules
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
