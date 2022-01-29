package app

import "fmt"

type Message struct {
	s string
}

func NewMessageFromStringer(s fmt.Stringer) Message {
	return Message{s: s.String()}
}

func NewMessageFromError(err error) Message {
	return Message{s: err.Error()}
}

func (m Message) String() string {
	return m.s
}
