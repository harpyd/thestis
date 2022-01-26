package mock

import (
	"context"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/specification"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

var (
	errNoSuchID                 = errors.New("no such id in mock map")
	errNoSpecWithTestCampaignID = errors.New("no specification with test campaign id in mock map")
)

type TestCampaignsRepository struct {
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
	tc, ok := m.campaigns[tcID]
	if !ok {
		return nil, app.NewTestCampaignNotFoundError(errNoSuchID)
	}

	return &tc, nil
}

func (m *TestCampaignsRepository) AddTestCampaign(_ context.Context, tc *testcampaign.TestCampaign) error {
	m.campaigns[tc.ID()] = *tc

	return nil
}

func (m *TestCampaignsRepository) UpdateTestCampaign(
	ctx context.Context,
	tcID string,
	updateFn app.TestCampaignUpdater,
) error {
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
	return len(m.campaigns)
}

type SpecificationsRepository struct {
	specifications map[string]specification.Specification
}

func NewSpecificationsRepository(specs ...*specification.Specification) *SpecificationsRepository {
	sm := &SpecificationsRepository{
		specifications: make(map[string]specification.Specification, len(specs)),
	}

	for _, spec := range specs {
		sm.specifications[spec.ID()] = *spec
	}

	return sm
}

func (m *SpecificationsRepository) GetSpecification(
	_ context.Context,
	specID string,
) (*specification.Specification, error) {
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
	for _, spec := range m.specifications {
		if spec.TestCampaignID() == tcID {
			return &spec, nil
		}
	}

	return nil, app.NewSpecificationNotFoundError(errNoSpecWithTestCampaignID)
}

func (m *SpecificationsRepository) AddSpecification(_ context.Context, spec *specification.Specification) error {
	m.specifications[spec.ID()] = *spec

	return nil
}

func (m *SpecificationsRepository) SpecificationsNumber() int {
	return len(m.specifications)
}
