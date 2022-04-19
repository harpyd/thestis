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

func newFlowDocument(flow *flow.Flow) flowDocument {
	return flowDocument{
		ID:            flow.ID(),
		PerformanceID: flow.PerformanceID(),
		Statuses:      newStatusDocuments(flow.Statuses()),
	}
}

func newStatusDocuments(statuses []*flow.Status) []statusDocument {
	documents := make([]statusDocument, 0, len(statuses))
	for _, s := range statuses {
		documents = append(documents, newStatusDocument(s))
	}

	return documents
}

func newStatusDocument(status *flow.Status) statusDocument {
	return statusDocument{
		Slug:           newScenarioSlugDocument(status.Slug()),
		State:          status.State(),
		ThesisStatuses: newThesisStatusDocuments(status.ThesisStatuses()),
	}
}

func newScenarioSlugDocument(slug specification.Slug) scenarioSlugDocument {
	return scenarioSlugDocument{
		Story:    slug.Story(),
		Scenario: slug.Scenario(),
	}
}

func newThesisStatusDocuments(statuses []*flow.ThesisStatus) []thesisStatusDocument {
	documents := make([]thesisStatusDocument, 0, len(statuses))
	for _, s := range statuses {
		documents = append(documents, newThesisStatusDocument(s))
	}

	return documents
}

func newThesisStatusDocument(status *flow.ThesisStatus) thesisStatusDocument {
	return thesisStatusDocument{
		ThesisSlug:   status.ThesisSlug(),
		State:        status.State(),
		OccurredErrs: status.OccurredErrs(),
	}
}

func (d flowDocument) toFlow() *flow.Flow {
	return flow.Unmarshal(flow.Params{
		ID:            d.ID,
		PerformanceID: d.PerformanceID,
		Statuses:      d.Statuses.toStatuses(),
	})
}

func (ds statusDocuments) toStatuses() []*flow.Status {
	transitions := make([]*flow.Status, 0, len(ds))
	for _, d := range ds {
		transitions = append(transitions, d.toStatus())
	}

	return transitions
}

func (d statusDocument) toStatus() *flow.Status {
	return flow.NewStatus(
		d.Slug.toSlug(),
		d.State,
		d.ThesisStatuses.toThesisStatuses()...,
	)
}

func (d scenarioSlugDocument) toSlug() specification.Slug {
	return specification.NewScenarioSlug(d.Story, d.Scenario)
}

func (ds thesisStatusDocuments) toThesisStatuses() []*flow.ThesisStatus {
	statuses := make([]*flow.ThesisStatus, 0, len(ds))
	for _, d := range ds {
		statuses = append(statuses, d.toThesisStatus())
	}

	return statuses
}

func (d thesisStatusDocument) toThesisStatus() *flow.ThesisStatus {
	return flow.NewThesisStatus(
		d.ThesisSlug,
		d.State,
		d.OccurredErrs...,
	)
}
