package mock

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/flow"
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

var errDuplicateID = errors.New("duplicate id in mock map")

type TestCampaignsRepository struct {
	mu        sync.RWMutex
	campaigns map[string]testcampaign.TestCampaign
}

func NewTestCampaignsRepository(tcs ...*testcampaign.TestCampaign) *TestCampaignsRepository {
	tcm := &TestCampaignsRepository{
		campaigns: make(map[string]testcampaign.TestCampaign, len(tcs)),
	}

	for _, tc := range tcs {
		tcm.campaigns[tc.ID()] = *tc
	}

	return tcm
}

func (m *TestCampaignsRepository) GetTestCampaign(ctx context.Context, tcID string) (*testcampaign.TestCampaign, error) {
	if ctx.Err() != nil {
		return nil, app.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	tc, ok := m.campaigns[tcID]
	if !ok {
		return nil, app.ErrTestCampaignNotFound
	}

	return &tc, nil
}

func (m *TestCampaignsRepository) AddTestCampaign(ctx context.Context, tc *testcampaign.TestCampaign) error {
	if ctx.Err() != nil {
		return app.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.campaigns[tc.ID()]; ok {
		return app.WrapWithDatabaseError(errDuplicateID)
	}

	m.campaigns[tc.ID()] = *tc

	return nil
}

func (m *TestCampaignsRepository) UpdateTestCampaign(
	ctx context.Context,
	tcID string,
	updateFn app.TestCampaignUpdater,
) error {
	if ctx.Err() != nil {
		return app.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	tc, ok := m.campaigns[tcID]
	if !ok {
		return app.ErrTestCampaignNotFound
	}

	updatedTC, err := updateFn(ctx, &tc)
	if err != nil {
		return err
	}

	m.campaigns[updatedTC.ID()] = *updatedTC

	return nil
}

func (m *TestCampaignsRepository) TestCampaignsNumber() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.campaigns)
}

type SpecificationsRepository struct {
	mu             sync.RWMutex
	specifications map[string]specification.Specification
}

func NewSpecificationsRepository(specs ...*specification.Specification) *SpecificationsRepository {
	m := &SpecificationsRepository{
		specifications: make(map[string]specification.Specification, len(specs)),
	}

	for _, spec := range specs {
		m.specifications[spec.ID()] = *spec
	}

	return m
}

func (m *SpecificationsRepository) GetSpecification(
	ctx context.Context,
	specID string,
) (*specification.Specification, error) {
	if ctx.Err() != nil {
		return nil, app.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	spec, ok := m.specifications[specID]
	if !ok {
		return nil, app.ErrSpecificationNotFound
	}

	return &spec, nil
}

func (m *SpecificationsRepository) GetActiveSpecificationByTestCampaignID(
	ctx context.Context,
	tcID string,
) (*specification.Specification, error) {
	if ctx.Err() != nil {
		return nil, app.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, spec := range m.specifications {
		if spec.TestCampaignID() == tcID {
			return &spec, nil
		}
	}

	return nil, app.ErrSpecificationNotFound
}

func (m *SpecificationsRepository) AddSpecification(
	ctx context.Context,
	spec *specification.Specification,
) error {
	if ctx.Err() != nil {
		return app.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.specifications[spec.ID()]; ok {
		return app.WrapWithDatabaseError(errDuplicateID)
	}

	m.specifications[spec.ID()] = *spec

	return nil
}

func (m *SpecificationsRepository) SpecificationsNumber() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.specifications)
}

type (
	PerformancesRepository struct {
		mu           sync.RWMutex
		performances map[string]performance.Performance
	}
)

func NewPerformancesRepository(perfs ...*performance.Performance) *PerformancesRepository {
	m := &PerformancesRepository{
		performances: make(map[string]performance.Performance, len(perfs)),
	}

	for _, p := range perfs {
		m.performances[p.ID()] = *p
	}

	return m
}

func (m *PerformancesRepository) GetPerformance(
	ctx context.Context,
	perfID string,
	_ app.SpecificationGetter,
	_ ...performance.Option,
) (*performance.Performance, error) {
	if ctx.Err() != nil {
		return nil, app.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	perf, ok := m.performances[perfID]
	if !ok {
		return nil, app.ErrPerformanceNotFound
	}

	return &perf, nil
}

func (m *PerformancesRepository) AddPerformance(ctx context.Context, perf *performance.Performance) error {
	if ctx.Err() != nil {
		return app.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.performances[perf.ID()]; ok {
		return app.WrapWithDatabaseError(errDuplicateID)
	}

	m.performances[perf.ID()] = *perf

	return nil
}

func (m *PerformancesRepository) PerformancesNumber() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.performances)
}

type FlowsRepository struct {
	mu    sync.RWMutex
	flows map[string]flow.Flow
}

func NewFlowsRepository(flows ...flow.Flow) *FlowsRepository {
	m := &FlowsRepository{
		flows: make(map[string]flow.Flow, len(flows)),
	}

	for _, f := range flows {
		m.flows[f.ID()] = f
	}

	return m
}

func (m *FlowsRepository) GetFlow(ctx context.Context, flowID string) (flow.Flow, error) {
	if ctx.Err() != nil {
		return flow.Flow{}, app.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	f, ok := m.flows[flowID]
	if !ok {
		return flow.Flow{}, app.ErrFlowNotFound
	}

	return f, nil
}

func (m *FlowsRepository) UpsertFlow(ctx context.Context, flow flow.Flow) error {
	if ctx.Err() != nil {
		return app.WrapWithDatabaseError(ctx.Err())
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.flows[flow.ID()] = flow

	return nil
}
