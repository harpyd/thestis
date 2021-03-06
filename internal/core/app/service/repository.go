package service

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/entity/flow"
	"github.com/harpyd/thestis/internal/core/entity/pipeline"
	"github.com/harpyd/thestis/internal/core/entity/specification"
	"github.com/harpyd/thestis/internal/core/entity/testcampaign"
)

var (
	ErrTestCampaignNotFound  = errors.New("test campaign not found")
	ErrSpecificationNotFound = errors.New("specification not found")
	ErrPipelineNotFound      = errors.New("pipeline not found")
	ErrFlowNotFound          = errors.New("flow not found")
)

type (
	TestCampaignRepository interface {
		GetTestCampaign(ctx context.Context, tcID string) (*testcampaign.TestCampaign, error)
		AddTestCampaign(ctx context.Context, tc *testcampaign.TestCampaign) error
		UpdateTestCampaign(ctx context.Context, tcID string, updater TestCampaignUpdater) error
	}

	TestCampaignUpdater func(
		ctx context.Context,
		tc *testcampaign.TestCampaign,
	) (*testcampaign.TestCampaign, error)
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
)

type (
	PipelineRepository interface {
		GetPipeline(
			ctx context.Context,
			pipeID string,
			specGetter SpecificationGetter,
			registrars ...pipeline.ExecutorRegistrar,
		) (*pipeline.Pipeline, error)
		AddPipeline(ctx context.Context, pipe *pipeline.Pipeline) error
	}

	SpecificationGetter interface {
		GetSpecification(ctx context.Context, specID string) (*specification.Specification, error)
	}
)

func AvailableSpecification(spec *specification.Specification) SpecificationGetter {
	return getSpecificationFunc(func() *specification.Specification {
		return spec
	})
}

func WithoutSpecification() SpecificationGetter {
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
