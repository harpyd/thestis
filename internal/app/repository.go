package app

import (
	"context"

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
