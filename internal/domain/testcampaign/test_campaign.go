package testcampaign

import (
	"time"

	"github.com/pkg/errors"
)

type TestCampaign struct {
	id       string
	viewName string
	summary  string

	activeSpecificationID string
	ownerID               string
	createdAt             time.Time
}

type Params struct {
	ID                    string
	ViewName              string
	Summary               string
	ActiveSpecificationID string
	OwnerID               string
	CreatedAt             time.Time
}

func New(params Params) (*TestCampaign, error) {
	if params.ID == "" {
		return nil, NewEmptyIDError()
	}

	if params.OwnerID == "" {
		return nil, NewEmptyOwnerIDError()
	}

	return &TestCampaign{
		id:                    params.ID,
		viewName:              params.ViewName,
		summary:               params.Summary,
		activeSpecificationID: params.ActiveSpecificationID,
		ownerID:               params.OwnerID,
		createdAt:             params.CreatedAt,
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

func (tc *TestCampaign) OwnerID() string {
	return tc.ownerID
}

func (tc *TestCampaign) CreatedAt() time.Time {
	return tc.createdAt
}

func (tc *TestCampaign) BindActiveSpecification(specificationID string) {
	tc.activeSpecificationID = specificationID
}

var (
	errEmptyID      = errors.New("empty test campaign ID")
	errEmptyOwnerID = errors.New("empty owner ID")
)

func NewEmptyIDError() error {
	return errEmptyID
}

func IsEmptyIDError(err error) bool {
	return errors.Is(err, errEmptyID)
}

func NewEmptyOwnerIDError() error {
	return errEmptyOwnerID
}

func IsEmptyOwnerIDError(err error) bool {
	return errors.Is(err, errEmptyOwnerID)
}
