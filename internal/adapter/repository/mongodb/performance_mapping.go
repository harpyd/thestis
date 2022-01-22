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
		From   string         `bson:"from"`
		To     string         `bson:"to"`
		Thesis thesisDocument `bson:"thesis"`
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
			From:   a.From(),
			To:     a.To(),
			Thesis: marshalToThesisDocument(a.Thesis()),
		})
	}

	return documents
}

func (d performanceDocument) unmarshalToPerformance() *performance.Performance {
	actions := make([]performance.ActionParam, 0, len(d.Actions))
	for _, a := range d.Actions {
		actions = append(actions, a.unmarshalToActionParam())
	}

	return performance.UnmarshalFromDatabase(performance.Params{
		ID:              d.ID,
		OwnerID:         d.OwnerID,
		SpecificationID: d.SpecificationID,
		Actions:         actions,
	})
}

func (d actionDocument) unmarshalToActionParam() performance.ActionParam {
	thesisBuilder := specification.NewThesisBuilder()
	buildFn := d.Thesis.unmarshalToThesisBuildFn()
	buildFn(thesisBuilder)

	return performance.ActionParam{
		From:   d.From,
		To:     d.To,
		Thesis: thesisBuilder.ErrlessBuild(d.Thesis.Slug),
	}
}
