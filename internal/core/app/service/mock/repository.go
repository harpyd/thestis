package mock

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/flow"
	"github.com/harpyd/thestis/internal/core/entity/performance"
	"github.com/harpyd/thestis/internal/core/entity/specification"
	"github.com/harpyd/thestis/internal/core/entity/testcampaign"
)

var errDuplicateID = errors.New("duplicate id in mock map")

type TestCampaignRepository struct {
	mu        sync.RWMutex
	campaigns map[string]testcampaign.TestCampaign
}

func NewTestCampaignRepository(tcs ...*testcampaign.TestCampaign) *TestCampaignRepository {
	tcm := &TestCampaignRepository{
		campaigns: make(map[string]testcampaign.TestCampaign, len(tcs)),
	}

	for _, tc := range tcs {
		tcm.campaigns[tc.ID()] = *tc
	}

	return tcm
}

func (m *TestCampaignRepository) GetTestCampaign(ctx context.Context, tcID string) (*testcampaign.TestCampaign, error) {
	if ctx.Err() != nil {
		return nil, service.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	tc, ok := m.campaigns[tcID]
	if !ok {
		return nil, service.ErrTestCampaignNotFound
	}

	return &tc, nil
}

func (m *TestCampaignRepository) AddTestCampaign(ctx context.Context, tc *testcampaign.TestCampaign) error {
	if ctx.Err() != nil {
		return service.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.campaigns[tc.ID()]; ok {
		return service.WrapWithDatabaseError(errDuplicateID)
	}

	m.campaigns[tc.ID()] = *tc

	return nil
}

func (m *TestCampaignRepository) UpdateTestCampaign(
	ctx context.Context,
	tcID string,
	updateFn service.TestCampaignUpdater,
) error {
	if ctx.Err() != nil {
		return service.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	tc, ok := m.campaigns[tcID]
	if !ok {
		return service.ErrTestCampaignNotFound
	}

	updatedTC, err := updateFn(ctx, &tc)
	if err != nil {
		return err
	}

	m.campaigns[updatedTC.ID()] = *updatedTC

	return nil
}

func (m *TestCampaignRepository) TestCampaignsNumber() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.campaigns)
}

type SpecificationRepository struct {
	mu             sync.RWMutex
	specifications map[string]specification.Specification
}

func NewSpecificationRepository(specs ...*specification.Specification) *SpecificationRepository {
	m := &SpecificationRepository{
		specifications: make(map[string]specification.Specification, len(specs)),
	}

	for _, spec := range specs {
		m.specifications[spec.ID()] = *spec
	}

	return m
}

func (m *SpecificationRepository) GetSpecification(
	ctx context.Context,
	specID string,
) (*specification.Specification, error) {
	if ctx.Err() != nil {
		return nil, service.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	spec, ok := m.specifications[specID]
	if !ok {
		return nil, service.ErrSpecificationNotFound
	}

	return &spec, nil
}

func (m *SpecificationRepository) GetActiveSpecificationByTestCampaignID(
	ctx context.Context,
	tcID string,
) (*specification.Specification, error) {
	if ctx.Err() != nil {
		return nil, service.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, spec := range m.specifications {
		if spec.TestCampaignID() == tcID {
			return &spec, nil
		}
	}

	return nil, service.ErrSpecificationNotFound
}

func (m *SpecificationRepository) AddSpecification(
	ctx context.Context,
	spec *specification.Specification,
) error {
	if ctx.Err() != nil {
		return service.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.specifications[spec.ID()]; ok {
		return service.WrapWithDatabaseError(errDuplicateID)
	}

	m.specifications[spec.ID()] = *spec

	return nil
}

func (m *SpecificationRepository) SpecificationsNumber() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.specifications)
}

type (
	PerformanceRepository struct {
		mu           sync.RWMutex
		performances map[string]performance.Performance
	}
)

func NewPerformanceRepository(perfs ...*performance.Performance) *PerformanceRepository {
	m := &PerformanceRepository{
		performances: make(map[string]performance.Performance, len(perfs)),
	}

	for _, p := range perfs {
		m.performances[p.ID()] = *p
	}

	return m
}

func (m *PerformanceRepository) GetPerformance(
	ctx context.Context,
	perfID string,
	_ service.SpecificationGetter,
	_ ...performance.Option,
) (*performance.Performance, error) {
	if ctx.Err() != nil {
		return nil, service.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	perf, ok := m.performances[perfID]
	if !ok {
		return nil, service.ErrPerformanceNotFound
	}

	return &perf, nil
}

func (m *PerformanceRepository) AddPerformance(ctx context.Context, perf *performance.Performance) error {
	if ctx.Err() != nil {
		return service.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.performances[perf.ID()]; ok {
		return service.WrapWithDatabaseError(errDuplicateID)
	}

	m.performances[perf.ID()] = *perf

	return nil
}

func (m *PerformanceRepository) PerformancesNumber() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.performances)
}

type FlowRepository struct {
	mu    sync.RWMutex
	flows map[string]flow.Flow
}

func NewFlowRepository(flows ...flow.Flow) *FlowRepository {
	m := &FlowRepository{
		flows: make(map[string]flow.Flow, len(flows)),
	}

	for _, f := range flows {
		m.flows[f.ID()] = f
	}

	return m
}

func (m *FlowRepository) GetFlow(ctx context.Context, flowID string) (*flow.Flow, error) {
	if ctx.Err() != nil {
		return nil, service.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	f, ok := m.flows[flowID]
	if !ok {
		return nil, service.ErrFlowNotFound
	}

	return &f, nil
}

func (m *FlowRepository) UpsertFlow(ctx context.Context, flow *flow.Flow) error {
	if ctx.Err() != nil {
		return service.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.flows[flow.ID()] = *flow

	return nil
}
