package service

import "github.com/harpyd/thestis/internal/core/entity/performance"

type MessageReactor func(msg Message)

type Message struct {
	s     string
	event performance.Event
	err   error
}

func NewMessageFromStep(s performance.Step) Message {
	return Message{
		s:     s.String(),
		event: s.Event(),
		err:   s.Err(),
	}
}

func NewMessageFromError(err error) Message {
	if err == nil {
		return Message{}
	}

	return Message{
		s:     err.Error(),
		event: performance.NoEvent,
		err:   err,
	}
}

func (m Message) String() string {
	return m.s
}

func (m Message) Err() error {
	return m.err
}

func (m Message) Event() performance.Event {
	return m.event
}
