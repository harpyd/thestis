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

func marshalToPerformanceDocument(perf *performance.Performance) performanceDocument {
	return performanceDocument{
		ID:              perf.ID(),
		OwnerID:         perf.OwnerID(),
		SpecificationID: perf.SpecificationID(),
		Started:         perf.Started(),
	}
}

func marshalToSlugDocument(slug specification.Slug) slugDocument {
	return slugDocument{
		Kind:     slug.Kind(),
		Story:    slug.Story(),
		Scenario: slug.Scenario(),
		Thesis:   slug.Thesis(),
	}
}

func (d performanceDocument) unmarshalToPerformance(
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
