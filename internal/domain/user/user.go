package user

import (
	"fmt"
	"github.com/harpyd/thestis/internal/domain/performance"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/domain/specification"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

func CanSeeTestCampaign(userID string, tc *testcampaign.TestCampaign) error {
	if userID == tc.OwnerID() {
		return nil
	}

	return NewCantSeeTestCampaignError(userID, tc.OwnerID())
}

func CanSeeSpecification(userID string, spec *specification.Specification) error {
	if userID == spec.OwnerID() {
		return nil
	}

	return NewCantSeeSpecificationError(userID, spec.OwnerID())
}

func CanSeePerformance(userID string, perf *performance.Performance) error {
	if userID == perf.OwnerID() {
		return nil
	}

	return NewCantSeePerformanceError(userID, perf.OwnerID())
}

type (
	cantSeeTestCampaignError struct {
		userID  string
		ownerID string
	}

	cantSeeSpecificationError struct {
		userID  string
		ownerID string
	}

	cantSeePerformanceError struct {
		userID  string
		ownerID string
	}
)

func NewCantSeeTestCampaignError(userID, ownerID string) error {
	return errors.WithStack(cantSeeTestCampaignError{
		userID:  userID,
		ownerID: ownerID,
	})
}

func IsCantSeeTestCampaignError(err error) bool {
	var target cantSeeTestCampaignError

	return errors.As(err, &target)
}

func (e cantSeeTestCampaignError) Error() string {
	return fmt.Sprintf("user %s can't see user %s test campaign", e.userID, e.ownerID)
}

func NewCantSeeSpecificationError(userID, ownerID string) error {
	return errors.WithStack(cantSeeSpecificationError{
		userID:  userID,
		ownerID: ownerID,
	})
}

func IsCantSeeSpecificationError(err error) bool {
	var target cantSeeSpecificationError

	return errors.As(err, &target)
}

func (e cantSeeSpecificationError) Error() string {
	return fmt.Sprintf("user %s can't see user %s specification", e.userID, e.ownerID)
}

func NewCantSeePerformanceError(userID, ownerID string) error {
	return errors.WithStack(cantSeePerformanceError{
		userID:  userID,
		ownerID: ownerID,
	})
}

func IsCantSeePerformanceError(err error) bool {
	var target cantSeePerformanceError

	return errors.As(err, &target)
}

func (e cantSeePerformanceError) Error() string {
	return fmt.Sprintf("user %s can't see user %s performance", e.userID, e.ownerID)
}
