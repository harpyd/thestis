package mongodb

import (
	"github.com/harpyd/thestis/internal/domain/flow"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type (
	flowDocument struct {
		ID            string          `bson:"_id"`
		OverallState  flow.State      `bson:"overallState"`
		PerformanceID string          `bson:"performanceId"`
		Statuses      statusDocuments `bson:"statuses"`
	}

	statusDocuments []statusDocument

	statusDocument struct {
		Slug           scenarioSlugDocument  `bson:"slug"`
		State          flow.State            `bson:"state"`
		ThesisStatuses thesisStatusDocuments `bson:"thesisStatuses"`
	}

	thesisStatusDocuments []thesisStatusDocument

	thesisStatusDocument struct {
		ThesisSlug   string     `bson:"thesisSlug"`
		State        flow.State `bson:"state"`
		OccurredErrs []string   `bson:"occurredErrs"`
	}

	scenarioSlugDocument struct {
		Story    string `bson:"story"`
		Scenario string `bson:"scenario"`
	}
)

func marshalToFlowDocument(flow flow.Flow) flowDocument {
	return flowDocument{
		ID:            flow.ID(),
		PerformanceID: flow.PerformanceID(),
		Statuses:      marshalToStatusDocuments(flow.Statuses()),
	}
}

func marshalToStatusDocuments(statuses []*flow.Status) []statusDocument {
	documents := make([]statusDocument, 0, len(statuses))
	for _, s := range statuses {
		documents = append(documents, marshalToStatusDocument(s))
	}

	return documents
}

func marshalToStatusDocument(status *flow.Status) statusDocument {
	return statusDocument{
		Slug:           marshalToScenarioSlugDocument(status.Slug()),
		State:          status.State(),
		ThesisStatuses: marshalToThesisStatusDocuments(status.ThesisStatuses()),
	}
}

func marshalToScenarioSlugDocument(slug specification.Slug) scenarioSlugDocument {
	return scenarioSlugDocument{
		Story:    slug.Story(),
		Scenario: slug.Scenario(),
	}
}

func marshalToThesisStatusDocuments(statuses []*flow.ThesisStatus) []thesisStatusDocument {
	documents := make([]thesisStatusDocument, 0, len(statuses))
	for _, s := range statuses {
		documents = append(documents, marshalToThesisStatusDocument(s))
	}

	return documents
}

func marshalToThesisStatusDocument(status *flow.ThesisStatus) thesisStatusDocument {
	return thesisStatusDocument{
		ThesisSlug:   status.ThesisSlug(),
		State:        status.State(),
		OccurredErrs: status.OccurredErrs(),
	}
}

func (d flowDocument) unmarshalToFlow() flow.Flow {
	return flow.Unmarshal(flow.Params{
		ID:            d.ID,
		PerformanceID: d.PerformanceID,
		Statuses:      d.Statuses.unmarshalToStatuses(),
	})
}

func (ds statusDocuments) unmarshalToStatuses() []*flow.Status {
	transitions := make([]*flow.Status, 0, len(ds))
	for _, d := range ds {
		transitions = append(transitions, d.unmarshalToStatus())
	}

	return transitions
}

func (d statusDocument) unmarshalToStatus() *flow.Status {
	return flow.NewStatus(
		d.Slug.unmarshalToSlug(),
		d.State,
		d.ThesisStatuses.unmarshalToThesisStatuses()...,
	)
}

func (d scenarioSlugDocument) unmarshalToSlug() specification.Slug {
	return specification.NewScenarioSlug(d.Story, d.Scenario)
}

func (ds thesisStatusDocuments) unmarshalToThesisStatuses() []*flow.ThesisStatus {
	statuses := make([]*flow.ThesisStatus, 0, len(ds))
	for _, d := range ds {
		statuses = append(statuses, d.unmarshalToThesisStatus())
	}

	return statuses
}

func (d thesisStatusDocument) unmarshalToThesisStatus() *flow.ThesisStatus {
	return flow.NewThesisStatus(
		d.ThesisSlug,
		d.State,
		d.OccurredErrs...,
	)
}
