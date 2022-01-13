package performance

import "fmt"

type Event interface {
	Err() error
	Performed() bool
	String() string
}

type actionEvent struct {
	from          string
	to            string
	performerType performerType
	performed     bool
}

func (e actionEvent) Err() error {
	return nil
}

func (e actionEvent) Performed() bool {
	return e.performed
}

func (e actionEvent) String() string {
	msg := fmt.Sprintf("Performance event: `%s -(%s)-> %s`", e.from, e.performerType, e.to)

	if e.performed {
		return fmt.Sprintf("%s performed", msg)
	}

	return fmt.Sprintf("%s not performed", msg)
}

type errEvent struct {
	err error
}

func (e errEvent) Err() error {
	return e.err
}

func (e errEvent) Performed() bool {
	return false
}

func (e errEvent) String() string {
	return fmt.Sprintf("Performance event: %s", e.err)
}
