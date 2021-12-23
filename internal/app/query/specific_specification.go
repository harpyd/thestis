package query

import (
	"context"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/app"
)

type SpecificSpecificationHandler struct {
	readModel app.SpecificSpecificationReadModel
}

func NewSpecificSpecificationHandler(readModel app.SpecificSpecificationReadModel) SpecificSpecificationHandler {
	if readModel == nil {
		panic("specific specification read model is nil")
	}

	return SpecificSpecificationHandler{
		readModel: readModel,
	}
}

func (h SpecificSpecificationHandler) Handle(
	ctx context.Context,
	qry app.SpecificSpecificationQuery,
) (app.SpecificSpecification, error) {
	specs, err := h.readModel.FindSpecification(ctx, qry)

	return specs, errors.Wrap(err, "getting specification")
}
