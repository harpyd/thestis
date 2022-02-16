package mongodb

import (
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type (
	performanceDocument struct {
		ID              string           `bson:"_id"`
		OwnerID         string           `bson:"ownerId"`
		SpecificationID string           `bson:"specificationId"`
		Actions         []actionDocument `bson:"actions"`
	}

	actionDocument struct {
		From          string                    `bson:"from"`
		To            string                    `bson:"to"`
		Thesis        thesisDocument            `bson:"thesis"`
		PerformerType performance.PerformerType `bson:"performerType"`
	}
)

func marshalToPerformanceDocument(perf *performance.Performance) performanceDocument {
	return performanceDocument{
		ID:              perf.ID(),
		OwnerID:         perf.OwnerID(),
		SpecificationID: perf.SpecificationID(),
		Actions:         marshalToActionDocuments(perf.Actions()),
	}
}

func marshalToActionDocuments(actions []performance.Action) []actionDocument {
	documents := make([]actionDocument, 0, len(actions))
	for _, a := range actions {
		documents = append(documents, actionDocument{
			From:          a.From(),
			To:            a.To(),
			Thesis:        marshalToThesisDocument(a.Thesis()),
			PerformerType: a.PerformerType(),
		})
	}

	return documents
}

func (d performanceDocument) unmarshalToPerformance() *performance.Performance {
	actions := make([]performance.Action, 0, len(d.Actions))
	for _, a := range d.Actions {
		actions = append(actions, a.unmarshalToAction())
	}

	return performance.Unmarshal(performance.Params{
		OwnerID:         d.OwnerID,
		SpecificationID: d.SpecificationID,
		Actions:         actions,
	}, performance.WithID(d.ID))
}

func (d actionDocument) unmarshalToAction() performance.Action {
	thesisBuilder := specification.NewThesisBuilder()
	buildFn := d.Thesis.unmarshalToThesisBuildFn()
	buildFn(thesisBuilder)

	return performance.NewAction(d.From, d.To, thesisBuilder.ErrlessBuild(d.Thesis.Slug), d.PerformerType)
}
