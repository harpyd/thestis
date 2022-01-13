package performance

type Attempt struct {
	flow    *Flow
	context *Context
}

func newAttempt(perf *Performance) Attempt {
	return Attempt{
		context: newContext(),
		flow:    newFlow(perf.actionGraph),
	}
}

func (a Attempt) Flow() *Flow {
	return a.flow
}

func (a Attempt) Context() *Context {
	return a.context
}
