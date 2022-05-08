package query

import (
	"context"

	"github.com/pkg/errors"
)

type SpecificSpecification struct {
	SpecificationID string
	UserID          string
}

type SpecificSpecificationHandler interface {
	Handle(ctx context.Context, qry SpecificSpecification) (SpecificSpecificationView, error)
}

type SpecificSpecificationReadModel interface {
	FindSpecification(ctx context.Context, qry SpecificSpecification) (SpecificSpecificationView, error)
}

type specificSpecificationHandler struct {
	readModel SpecificSpecificationReadModel
}

func NewSpecificSpecificationHandler(readModel SpecificSpecificationReadModel) SpecificSpecificationHandler {
	if readModel == nil {
		panic("specific specification read model is nil")
	}

	return specificSpecificationHandler{
		readModel: readModel,
	}
}

func (h specificSpecificationHandler) Handle(
	ctx context.Context,
	qry SpecificSpecification,
) (SpecificSpecificationView, error) {
	specs, err := h.readModel.FindSpecification(ctx, qry)

	return specs, errors.Wrap(err, "getting specification")
}
