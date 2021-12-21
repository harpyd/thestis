package testcampaign

import "github.com/pkg/errors"

type TestCampaign struct {
	id       string
	viewName string
	summary  string

	activeSpecificationID string
	userID                string
}

type Params struct {
	ID                    string
	ViewName              string
	Summary               string
	ActiveSpecificationID string
	UserID                string
}

func New(params Params) (*TestCampaign, error) {
	if params.ID == "" {
		return nil, NewEmptyIDError()
	}

	if params.UserID == "" {
		return nil, NewEmptyUserIDError()
	}

	return &TestCampaign{
		id:                    params.ID,
		viewName:              params.ViewName,
		summary:               params.Summary,
		activeSpecificationID: params.ActiveSpecificationID,
		userID:                params.UserID,
	}, nil
}

func (tc *TestCampaign) ID() string {
	return tc.id
}

func (tc *TestCampaign) ViewName() string {
	return tc.viewName
}

func (tc *TestCampaign) Summary() string {
	return tc.summary
}

func (tc *TestCampaign) ActiveSpecificationID() string {
	return tc.activeSpecificationID
}

func (tc *TestCampaign) UserID() string {
	return tc.userID
}

func (tc *TestCampaign) SetActiveSpecificationID(specificationID string) {
	tc.activeSpecificationID = specificationID
}

var (
	errEmptyID     = errors.New("empty test campaign ID")
	errEmptyUserID = errors.New("empty user ID")
)

func NewEmptyIDError() error {
	return errEmptyID
}

func IsEmptyIDError(err error) bool {
	return errors.Is(err, errEmptyID)
}

func NewEmptyUserIDError() error {
	return errEmptyUserID
}

func IsEmptyUserIDError(err error) bool {
	return errors.Is(err, errEmptyUserID)
}
