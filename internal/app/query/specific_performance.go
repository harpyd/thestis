package query

import (
	"context"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/app"
)

type SpecificPerformanceHandler struct {
	readModel app.SpecificPerformanceReadModel
}

func NewSpecificPerformanceHandler(readModel app.SpecificPerformanceReadModel) SpecificPerformanceHandler {
	if readModel == nil {
		panic("specific performance read model is nil")
	}

	return SpecificPerformanceHandler{
		readModel: readModel,
	}
}

func (h SpecificPerformanceHandler) Handle(
	ctx context.Context,
	qry app.SpecificPerformanceQuery,
) (app.SpecificPerformance, error) {
	perf, err := h.readModel.FindPerformance(ctx, qry)

	return perf, errors.Wrap(err, "getting performance")
}
