package user

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/entity/pipeline"
	"github.com/harpyd/thestis/internal/core/entity/specification"
	"github.com/harpyd/thestis/internal/core/entity/testcampaign"
)

type Permission string

const (
	NoPermission Permission = ""
	Read         Permission = "read"
	Write        Permission = "write"
)

type Resource string

const (
	NoResource    Resource = ""
	TestCampaign  Resource = "test campaign"
	Specification Resource = "specification"
	Pipeline      Resource = "pipeline"
)

type AccessError struct {
	userID     string
	resourceID string
	resource   Resource
	permission Permission
}

func NewAccessError(
	userID,
	resourceID string,
	resource Resource,
	perm Permission,
) error {
	return errors.WithStack(&AccessError{
		userID:     userID,
		resourceID: resourceID,
		resource:   resource,
		permission: perm,
	})
}

func (e *AccessError) UserID() string {
	return e.userID
}

func (e *AccessError) ResourceID() string {
	return e.resourceID
}

func (e *AccessError) Resource() Resource {
	return e.resource
}

func (e *AccessError) Permission() Permission {
	return e.permission
}

func (e *AccessError) Error() string {
	if e == nil {
		return ""
	}

	var b strings.Builder

	if e.userID != "" {
		_, _ = fmt.Fprintf(&b, "user #%s ", e.userID)
	}

	_, _ = b.WriteString("can't access")

	if e.resource != NoResource {
		_, _ = b.WriteString(" " + string(e.resource))
	}

	if e.resourceID != "" {
		_, _ = fmt.Fprintf(&b, " #%s", e.resourceID)
	}

	if e.permission != NoPermission {
		_, _ = fmt.Fprintf(&b, " with %q permission", e.permission)
	}

	return b.String()
}

func CanAccessTestCampaign(
	userID string,
	tc *testcampaign.TestCampaign,
	perm Permission,
) error {
	if userID == tc.OwnerID() {
		return nil
	}

	return NewAccessError(userID, tc.OwnerID(), TestCampaign, perm)
}

func CanAccessSpecification(
	userID string,
	spec *specification.Specification,
	perm Permission,
) error {
	if userID == spec.OwnerID() {
		return nil
	}

	return NewAccessError(userID, spec.OwnerID(), Specification, perm)
}

func CanAccessPipeline(
	userID string,
	pipe *pipeline.Pipeline,
	perm Permission,
) error {
	if userID == pipe.OwnerID() {
		return nil
	}

	return NewAccessError(userID, pipe.OwnerID(), Pipeline, perm)
}
