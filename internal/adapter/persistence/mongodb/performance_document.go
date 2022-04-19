package mongodb

import (
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type performanceDocument struct {
	ID              string `bson:"_id"`
	OwnerID         string `bson:"ownerId"`
	SpecificationID string `bson:"specificationId"`
	Started         bool   `bson:"started"`
}

func newPerformanceDocument(perf *performance.Performance) performanceDocument {
	return performanceDocument{
		ID:              perf.ID(),
		OwnerID:         perf.OwnerID(),
		SpecificationID: perf.SpecificationID(),
		Started:         perf.Started(),
	}
}

func (d performanceDocument) toPerformance(
	spec *specification.Specification,
	opts []performance.Option,
) *performance.Performance {
	return performance.Unmarshal(performance.Params{
		ID:            d.ID,
		Specification: spec,
		OwnerID:       d.OwnerID,
		Started:       d.Started,
	}, opts...)
}