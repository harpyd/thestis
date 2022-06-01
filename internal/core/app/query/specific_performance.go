package query

import (
	"context"

	"github.com/pkg/errors"
)

type SpecificPerformance struct {
	PerformanceID string
	UserID        string
}

type SpecificPerformanceHandler interface {
	Handle(ctx context.Context, qry SpecificPerformance) (SpecificPerformanceModel, error)
}

type SpecificPerformanceReadModel interface {
	FindPerformance(ctx context.Context, qry SpecificPerformance) (SpecificPerformanceModel, error)
}

type specificPerformanceHandler struct {
	readModel SpecificPerformanceReadModel
}

func NewSpecificPerformanceHandler(readModel SpecificPerformanceReadModel) SpecificPerformanceHandler {
	if readModel == nil {
		panic("specific performance read model is nil")
	}

	return specificPerformanceHandler{
		readModel: readModel,
	}
}

func (h specificPerformanceHandler) Handle(
	ctx context.Context,
	qry SpecificPerformance,
) (SpecificPerformanceModel, error) {
	perf, err := h.readModel.FindPerformance(ctx, qry)

	return perf, errors.Wrap(err, "getting performance")
}
