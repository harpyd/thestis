package app

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

type (
	TestCampaignsRepository interface {
		GetTestCampaign(ctx context.Context, tcID string) (*testcampaign.TestCampaign, error)
		AddTestCampaign(ctx context.Context, tc *testcampaign.TestCampaign) error
		UpdateTestCampaign(ctx context.Context, tcID string, updateFn TestCampaignUpdater) error
	}

	TestCampaignUpdater func(
		ctx context.Context,
		tc *testcampaign.TestCampaign,
	) (*testcampaign.TestCampaign, error)

	SpecificTestCampaignReadModel interface {
		FindTestCampaign(ctx context.Context, qry SpecificTestCampaignQuery) (SpecificTestCampaign, error)
	}
)

type (
	SpecificationsRepository interface {
		GetSpecification(ctx context.Context, specID string) (*specification.Specification, error)
		GetActiveSpecificationByTestCampaignID(
			ctx context.Context,
			testCampaignID string,
		) (*specification.Specification, error)
		AddSpecification(ctx context.Context, spec *specification.Specification) error
	}

	SpecificSpecificationReadModel interface {
		FindSpecification(ctx context.Context, qry SpecificSpecificationQuery) (SpecificSpecification, error)
	}
)

type (
	PerformancesRepository interface {
		GetPerformance(ctx context.Context, perfID string) (*performance.Performance, error)
		AddPerformance(ctx context.Context, perf *performance.Performance) error
		ExclusivelyDoWithPerformance(
			ctx context.Context,
			perf *performance.Performance,
			action PerformanceAction,
		) error
	}

	PerformanceAction func(perf *performance.Performance)
)

type FlowsRepository interface {
	GetFlow(ctx context.Context, flowID string) (performance.Flow, error)
	UpsertFlow(ctx context.Context, flow performance.Flow) error
}

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

	flowNotFoundError struct {
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

func NewFlowNotFoundError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(flowNotFoundError{
		err: err,
	})
}

func IsFlowNotFoundError(err error) bool {
	var target flowNotFoundError

	return errors.As(err, &target)
}

func (e flowNotFoundError) Cause() error {
	return e.err
}

func (e flowNotFoundError) Unwrap() error {
	return e.err
}

func (e flowNotFoundError) Error() string {
	return fmt.Sprintf("flow not found: %s", e.err)
}
