package testcampaign

import (
	"time"

	"github.com/pkg/errors"
)

type TestCampaign struct {
	id       string
	viewName string
	summary  string

	ownerID   string
	createdAt time.Time
}

type Params struct {
	ID        string
	ViewName  string
	Summary   string
	OwnerID   string
	CreatedAt time.Time
}

func MustNew(params Params) *TestCampaign {
	tc, err := New(params)
	if err != nil {
		panic(err)
	}

	return tc
}

var (
	ErrEmptyID      = errors.New("empty test campaign ID")
	ErrEmptyOwnerID = errors.New("empty owner ID")
)

func New(params Params) (*TestCampaign, error) {
	if params.ID == "" {
		return nil, ErrEmptyID
	}

	if params.OwnerID == "" {
		return nil, ErrEmptyOwnerID
	}

	return &TestCampaign{
		id:        params.ID,
		viewName:  params.ViewName,
		summary:   params.Summary,
		ownerID:   params.OwnerID,
		createdAt: params.CreatedAt,
	}, nil
}

func (tc *TestCampaign) ID() string {
	return tc.id
}

func (tc *TestCampaign) ViewName() string {
	return tc.viewName
}

func (tc *TestCampaign) SetViewName(viewName string) {
	tc.viewName = viewName
}

func (tc *TestCampaign) Summary() string {
	return tc.summary
}

func (tc *TestCampaign) SetSummary(summary string) {
	tc.summary = summary
}

func (tc *TestCampaign) OwnerID() string {
	return tc.ownerID
}

func (tc *TestCampaign) CreatedAt() time.Time {
	return tc.createdAt
}
