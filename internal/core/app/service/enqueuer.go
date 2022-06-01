package service

type Enqueuer interface {
	Enqueue(fn func())
}

type EnqueueFunc func(fn func())

func (e EnqueueFunc) Enqueue(fn func()) {
	e(fn)
}
