package performance

import "fmt"

type Event interface {
	String() string
	Err() error
}

type performEvent struct {
	from          string
	to            string
	performerType performerType
}

func (e performEvent) String() string {
	return fmt.Sprintf("Performance event: `%s -(%s)-> %s`", e.from, e.performerType, e.to)
}

func (e performEvent) Err() error {
	return nil
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
