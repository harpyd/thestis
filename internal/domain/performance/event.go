package performance

import "fmt"

type Event interface {
	Err() error
	String() string
}

type actionEvent struct {
	from          string
	to            string
	performerType performerType
}

func (e actionEvent) Err() error {
	return nil
}

func (e actionEvent) String() string {
	return fmt.Sprintf("Performance event: `%s -(%s)-> %s`", e.from, e.performerType, e.to)
}

type errEvent struct {
	err error
}

func (e errEvent) Err() error {
	return e.err
}

func (e errEvent) String() string {
	return fmt.Sprintf("Performance event: %s", e.err)
}
