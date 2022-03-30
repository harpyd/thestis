package mock

type Enqueuer struct {
	enqueueCalls int
}

func NewEnqueuer() *Enqueuer {
	return &Enqueuer{}
}

func (e *Enqueuer) Enqueue(fn func()) {
	e.enqueueCalls++

	go fn()
}

func (e *Enqueuer) EnqueueCalls() int {
	return e.enqueueCalls
}
