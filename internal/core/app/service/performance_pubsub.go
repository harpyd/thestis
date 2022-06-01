package service

import (
	"fmt"

	"github.com/pkg/errors"
)

type (
	PerformanceCancelPublisher interface {
		PublishPerformanceCancel(perfID string) error
	}

	PerformanceCancelSubscriber interface {
		SubscribePerformanceCancel(perfID string) (<-chan CancelSignal, error)
	}

	CancelSignal = struct{}
)

type PublishCancelError struct {
	err error
}

func WrapWithPublishCancelError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(&PublishCancelError{err: err})
}

func (e *PublishCancelError) Unwrap() error {
	return e.err
}

func (e *PublishCancelError) Error() string {
	if e == nil || e.err == nil {
		return ""
	}

	return fmt.Sprintf("publish cancel: %s", e.err)
}

type SubscribeCancelError struct {
	err error
}

func WrapWithSubscribeCancelError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(&SubscribeCancelError{err: err})
}

func (e *SubscribeCancelError) Unwrap() error {
	return e.err
}

func (e *SubscribeCancelError) Error() string {
	if e == nil || e.err == nil {
		return ""
	}

	return fmt.Sprintf("subscribe cancel: %s", e.err)
}
