package performance

type Attempt struct {
	state   *State
	context *Context
}

func newAttempt() Attempt {
	return Attempt{
		context: newContext(),
	}
}

func (a Attempt) State() *State {
	return a.state
}

func (a Attempt) Context() *Context {
	return a.context
}
