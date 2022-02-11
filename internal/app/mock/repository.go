package mock

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

var (
	errNoSuchID                 = errors.New("no such id in mock map")
	errDuplicateID              = errors.New("duplicate id in mock map")
	errNoSpecWithTestCampaignID = errors.New("no specification with test campaign id in mock map")
)

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

func (m *TestCampaignsRepository) GetTestCampaign(_ context.Context, tcID string) (*testcampaign.TestCampaign, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tc, ok := m.campaigns[tcID]
	if !ok {
		return nil, app.NewTestCampaignNotFoundError(errNoSuchID)
	}

	return &tc, nil
}

func (m *TestCampaignsRepository) AddTestCampaign(_ context.Context, tc *testcampaign.TestCampaign) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.campaigns[tc.ID()]; ok {
		return app.NewAlreadyExistsError(errDuplicateID)
	}

	m.campaigns[tc.ID()] = *tc

	return nil
}

func (m *TestCampaignsRepository) UpdateTestCampaign(
	ctx context.Context,
	tcID string,
	updateFn app.TestCampaignUpdater,
) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	tc, ok := m.campaigns[tcID]
	if !ok {
		return app.NewTestCampaignNotFoundError(errNoSuchID)
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
	_ context.Context,
	specID string,
) (*specification.Specification, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	spec, ok := m.specifications[specID]
	if !ok {
		return nil, app.NewSpecificationNotFoundError(errNoSuchID)
	}

	return &spec, nil
}

func (m *SpecificationsRepository) GetActiveSpecificationByTestCampaignID(
	_ context.Context,
	tcID string,
) (*specification.Specification, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, spec := range m.specifications {
		if spec.TestCampaignID() == tcID {
			return &spec, nil
		}
	}

	return nil, app.NewSpecificationNotFoundError(errNoSpecWithTestCampaignID)
}

func (m *SpecificationsRepository) AddSpecification(_ context.Context, spec *specification.Specification) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.specifications[spec.ID()]; ok {
		return app.NewAlreadyExistsError(errDuplicateID)
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
		performances map[string]lockedPerformance
	}

	lockedPerformance struct {
		lock        uint32
		performance performance.Performance
	}
)

func NewPerformancesRepository(perfs ...*performance.Performance) *PerformancesRepository {
	m := &PerformancesRepository{
		performances: make(map[string]lockedPerformance, len(perfs)),
	}

	for _, p := range perfs {
		m.performances[p.ID()] = lockedPerformance{
			performance: *p,
		}
	}

	return m
}

func (m *PerformancesRepository) GetPerformance(_ context.Context, perfID string) (*performance.Performance, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	lp, ok := m.performances[perfID]
	if !ok {
		return nil, app.NewPerformanceNotFoundError(errNoSuchID)
	}

	return &lp.performance, nil
}

func (m *PerformancesRepository) AddPerformance(_ context.Context, perf *performance.Performance) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.performances[perf.ID()]; ok {
		return app.NewAlreadyExistsError(errDuplicateID)
	}

	m.performances[perf.ID()] = lockedPerformance{
		performance: *perf,
	}

	return nil
}

func (m *PerformancesRepository) ExclusivelyDoWithPerformance(
	ctx context.Context,
	perf *performance.Performance,
	action app.PerformanceAction,
) error {
	m.mu.RLock()

	lp, ok := m.performances[perf.ID()]
	if !ok {
		return app.NewPerformanceNotFoundError(errNoSuchID)
	}

	m.mu.RUnlock()

	if !atomic.CompareAndSwapUint32(&lp.lock, 0, 1) {
		return performance.NewAlreadyStartedError()
	}

	go func() {
		defer atomic.StoreUint32(&lp.lock, 0)

		action(ctx, perf)
	}()

	return nil
}

func (m *PerformancesRepository) PerformancesNumber() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.performances)
}

type FlowsRepository struct {
	mu    sync.RWMutex
	flows map[string]performance.Flow
}

func NewFlowsRepository(flows ...performance.Flow) *FlowsRepository {
	m := &FlowsRepository{
		flows: make(map[string]performance.Flow, len(flows)),
	}

	for _, f := range flows {
		m.flows[f.ID()] = f
	}

	return m
}

func (m *FlowsRepository) GetFlow(_ context.Context, flowID string) (performance.Flow, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	flow, ok := m.flows[flowID]
	if !ok {
		return performance.Flow{}, app.NewFlowNotFoundError(errNoSuchID)
	}

	return flow, nil
}

func (m *FlowsRepository) UpsertFlow(_ context.Context, flow performance.Flow) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.flows[flow.ID()] = flow

	return nil
}
