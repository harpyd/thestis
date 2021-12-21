package user

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/domain/specification"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

func CanSeeTestCampaign(userID string, tc *testcampaign.TestCampaign) error {
	if userID == tc.OwnerID() {
		return nil
	}

	return NewUserCantSeeTestCampaignError(userID, tc.OwnerID())
}

func CanSeeSpecification(userID string, spec *specification.Specification) error {
	if userID == spec.OwnerID() {
		return nil
	}

	return NewUserCantSeeSpecificationError(userID, spec.OwnerID())
}

type (
	userCantSeeTestCampaignError struct {
		userID  string
		ownerID string
	}

	userCantSeeSpecificationError struct {
		userID  string
		ownerID string
	}
)

func NewUserCantSeeTestCampaignError(userID, ownerID string) error {
	return errors.WithStack(userCantSeeTestCampaignError{
		userID:  userID,
		ownerID: ownerID,
	})
}

func IsUserCantSeeTestCampaignError(err error) bool {
	var target userCantSeeTestCampaignError

	return errors.As(err, &target)
}

func (e userCantSeeTestCampaignError) Error() string {
	return fmt.Sprintf("user %s can't see user %s test campaign", e.userID, e.ownerID)
}

func NewUserCantSeeSpecificationError(userID, ownerID string) error {
	return errors.WithStack(userCantSeeSpecificationError{
		userID:  userID,
		ownerID: ownerID,
	})
}

func IsUserCantSeeSpecificationError(err error) bool {
	var target userCantSeeSpecificationError

	return errors.As(err, &target)
}

func (e userCantSeeSpecificationError) Error() string {
	return fmt.Sprintf("user %s can't see user %s specification", e.userID, e.ownerID)
}
