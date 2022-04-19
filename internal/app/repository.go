package app

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/domain/flow"
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

var (
	ErrTestCampaignNotFound  = errors.New("test campaign not found")
	ErrSpecificationNotFound = errors.New("specification not found")
	ErrPerformanceNotFound   = errors.New("performance not found")
	ErrFlowNotFound          = errors.New("flow not found")
)

type (
	TestCampaignRepository interface {
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
	SpecificationRepository interface {
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
	PerformanceRepository interface {
		GetPerformance(
			ctx context.Context,
			perfID string,
			specGetter SpecificationGetter,
			opts ...performance.Option,
		) (*performance.Performance, error)
		AddPerformance(ctx context.Context, perf *performance.Performance) error
	}

	SpecificationGetter interface {
		GetSpecification(ctx context.Context, specID string) (*specification.Specification, error)
	}

	SpecificPerformanceReadModel interface {
		FindPerformance(ctx context.Context, qry SpecificPerformanceQuery) (SpecificPerformance, error)
	}
)

func AvailableSpecification(spec *specification.Specification) SpecificationGetter {
	return getSpecificationFunc(func() *specification.Specification {
		return spec
	})
}

func DontGetSpecification() SpecificationGetter {
	return getSpecificationFunc(func() *specification.Specification {
		return nil
	})
}

type getSpecificationFunc func() *specification.Specification

func (f getSpecificationFunc) GetSpecification(_ context.Context, _ string) (*specification.Specification, error) {
	return f(), nil
}

type FlowRepository interface {
	GetFlow(ctx context.Context, flowID string) (*flow.Flow, error)
	UpsertFlow(ctx context.Context, flow *flow.Flow) error
}

type DatabaseError struct {
	err error
}

func WrapWithDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(&DatabaseError{
		err: err,
	})
}

func (e *DatabaseError) Unwrap() error {
	return e.err
}

func (e *DatabaseError) Error() string {
	if e == nil || e.err == nil {
		return ""
	}

	return fmt.Sprintf("database problem: %s", e.err)
}
