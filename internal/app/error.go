package app

import (
	"fmt"

	"github.com/pkg/errors"
)

type (
	databaseError struct {
		err error
	}

	testCampaignNotFoundError struct {
		err error
	}

	specificationNotFoundError struct {
		err error
	}

	performanceNotFoundError struct {
		err error
	}
)

func NewDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(databaseError{
		err: err,
	})
}

func IsDatabaseError(err error) bool {
	var target databaseError

	return errors.As(err, &target)
}

func (e databaseError) Cause() error {
	return e.err
}

func (e databaseError) Unwrap() error {
	return e.err
}

func (e databaseError) Error() string {
	return fmt.Sprintf("database problem: %s", e.err)
}

func NewTestCampaignNotFoundError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(testCampaignNotFoundError{
		err: err,
	})
}

func IsTestCampaignNotFoundError(err error) bool {
	var target testCampaignNotFoundError

	return errors.As(err, &target)
}

func (e testCampaignNotFoundError) Cause() error {
	return e.err
}

func (e testCampaignNotFoundError) Unwrap() error {
	return e.err
}

func (e testCampaignNotFoundError) Error() string {
	return fmt.Sprintf("test campaign not found: %s", e.err)
}

func NewSpecificationNotFoundError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(specificationNotFoundError{
		err: err,
	})
}

func IsSpecificationNotFoundError(err error) bool {
	var target specificationNotFoundError

	return errors.As(err, &target)
}

func (e specificationNotFoundError) Cause() error {
	return e.err
}

func (e specificationNotFoundError) Unwrap() error {
	return e.err
}

func (e specificationNotFoundError) Error() string {
	return fmt.Sprintf("specification not found: %s", e.err)
}

func NewPerformanceNotFoundError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(performanceNotFoundError{
		err: err,
	})
}

func IsPerformanceNotFoundError(err error) bool {
	var target performanceNotFoundError

	return errors.As(err, &target)
}

func (e performanceNotFoundError) Cause() error {
	return e.err
}

func (e performanceNotFoundError) Unwrap() error {
	return e.err
}

func (e performanceNotFoundError) Error() string {
	return fmt.Sprintf("performance not found: %s", e.err)
}
