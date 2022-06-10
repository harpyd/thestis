package query

import (
	"context"

	"github.com/pkg/errors"
)

type Specification struct {
	SpecificationID string
	UserID          string
}

type SpecificationHandler interface {
	Handle(ctx context.Context, qry Specification) (SpecificationModel, error)
}

type SpecificationReadModel interface {
	FindSpecification(ctx context.Context, qry Specification) (SpecificationModel, error)
}

type specificationHandler struct {
	readModel SpecificationReadModel
}

func NewSpecificationHandler(readModel SpecificationReadModel) SpecificationHandler {
	if readModel == nil {
		panic("specification read model is nil")
	}

	return specificationHandler{
		readModel: readModel,
	}
}

func (h specificationHandler) Handle(
	ctx context.Context,
	qry Specification,
) (SpecificationModel, error) {
	specs, err := h.readModel.FindSpecification(ctx, qry)

	return specs, errors.Wrap(err, "getting specification")
}
