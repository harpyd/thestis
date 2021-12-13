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