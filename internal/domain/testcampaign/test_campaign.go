package testcampaign

import "github.com/pkg/errors"

type TestCampaign struct {
	id                    string
	viewName              string
	activeSpecificationID string
}

func New(id string, viewName string) (*TestCampaign, error) {
	if id == "" {
		return nil, NewEmptyIDError()
	}

	return &TestCampaign{
		id:       id,
		viewName: viewName,
	}, nil
}

func UnmarshalFromDatabase(
	id string,
	viewName string,
	activeSpecificationID string,
) *TestCampaign {
	return &TestCampaign{
		id:                    id,
		viewName:              viewName,
		activeSpecificationID: activeSpecificationID,
	}
}

func (tc *TestCampaign) ID() string {
	return tc.id
}

func (tc *TestCampaign) ViewName() string {
	return tc.viewName
}

func (tc *TestCampaign) ActiveSpecificationID() string {
	return tc.activeSpecificationID
}

func (tc *TestCampaign) SetActiveSpecificationID(specificationID string) {
	tc.activeSpecificationID = specificationID
}

var errEmptyID = errors.New("empty test campaign ID")

func NewEmptyIDError() error {
	return errEmptyID
}

func IsEmptyIDError(err error) bool {
	return errors.Is(err, errEmptyID)
}
