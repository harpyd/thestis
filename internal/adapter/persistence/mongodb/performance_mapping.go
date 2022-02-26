package mongodb

import (
	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type (
	performanceDocument struct {
		ID              string           `bson:"_id"`
		OwnerID         string           `bson:"ownerId"`
		SpecificationID string           `bson:"specificationId"`
		Actions         []actionDocument `bson:"actions"`
		Started         bool             `bson:"started"`
	}

	actionDocument struct {
		From          string                    `bson:"from"`
		To            string                    `bson:"to"`
		Thesis        performanceThesisDocument `bson:"thesis"`
		PerformerType performance.PerformerType `bson:"performerType"`
	}

	performanceThesisDocument struct {
		Slug      slugDocument      `bson:"slug"`
		After     []string          `bson:"after"`
		Statement statementDocument `bson:"statement"`
		HTTP      httpDocument      `bson:"http"`
		Assertion assertionDocument `bson:"assertion"`
	}

	slugDocument struct {
		Story    string
		Scenario string
		Thesis   string
	}
)

func marshalToPerformanceDocument(perf *performance.Performance) performanceDocument {
	return performanceDocument{
		ID:              perf.ID(),
		OwnerID:         perf.OwnerID(),
		SpecificationID: perf.SpecificationID(),
		Actions:         marshalToActionDocuments(perf.Actions()),
		Started:         perf.Started(),
	}
}

func marshalToActionDocuments(actions []performance.Action) []actionDocument {
	documents := make([]actionDocument, 0, len(actions))
	for _, a := range actions {
		documents = append(documents, actionDocument{
			From:          a.From(),
			To:            a.To(),
			Thesis:        marshalToPerformanceThesisDocument(a.Thesis()),
			PerformerType: a.PerformerType(),
		})
	}

	return documents
}

func marshalToPerformanceThesisDocument(thesis specification.Thesis) performanceThesisDocument {
	return performanceThesisDocument{
		Slug:  marshalToSlugDocument(thesis.Slug()),
		After: marshalSlugsToStrings(thesis.Dependencies()),
		Statement: statementDocument{
			Keyword:  thesis.Statement().Stage().String(),
			Behavior: thesis.Statement().Behavior(),
		},
		HTTP:      marshalToHTTPDocument(thesis.HTTP()),
		Assertion: marshalToAssertionDocument(thesis.Assertion()),
	}
}

func marshalToSlugDocument(slug specification.Slug) slugDocument {
	return slugDocument{
		Story:    slug.Story(),
		Scenario: slug.Scenario(),
		Thesis:   slug.Thesis(),
	}
}

func (d performanceDocument) unmarshalToPerformance(opts app.PerformerOptions) *performance.Performance {
	actions := make([]performance.Action, 0, len(d.Actions))
	for _, a := range d.Actions {
		actions = append(actions, a.unmarshalToAction())
	}

	options := append(opts.ToPerformanceOptions(), performance.WithID(d.ID))

	return performance.Unmarshal(performance.Params{
		OwnerID:         d.OwnerID,
		SpecificationID: d.SpecificationID,
		Actions:         actions,
		Started:         d.Started,
	}, options...)
}

func (d actionDocument) unmarshalToAction() performance.Action {
	return performance.NewAction(
		d.From,
		d.To,
		d.Thesis.unmarshalToThesis(),
		d.PerformerType,
	)
}

func (d performanceThesisDocument) unmarshalToThesis() specification.Thesis {
	b := specification.NewThesisBuilder().
		WithAssertion(d.Assertion.unmarshalToAssertionBuildFn()).
		WithStatement(d.Statement.Keyword, d.Statement.Behavior).
		WithHTTP(d.HTTP.unmarshalToHTTPBuildFn())

	for _, dep := range d.After {
		b.WithDependencies(dep)
	}

	return b.ErrlessBuild(
		specification.NewThesisSlug(d.Slug.Story, d.Slug.Scenario, d.Slug.Thesis),
	)
}
