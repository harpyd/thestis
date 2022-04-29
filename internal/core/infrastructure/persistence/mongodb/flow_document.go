package mongodb

import (
	flow2 "github.com/harpyd/thestis/internal/core/domain/flow"
	"github.com/harpyd/thestis/internal/core/domain/specification"
)

type (
	flowDocument struct {
		ID            string           `bson:"_id"`
		PerformanceID string           `bson:"performanceId"`
		OverallState  flow2.State      `bson:"overallState"`
		Statuses      []statusDocument `bson:"statuses"`
	}

	statusDocument struct {
		Slug           scenarioSlugDocument  `bson:"slug"`
		State          flow2.State           `bson:"state"`
		ThesisStatuses thesisStatusDocuments `bson:"thesisStatuses"`
	}

	thesisStatusDocuments []thesisStatusDocument

	thesisStatusDocument struct {
		ThesisSlug   string      `bson:"thesisSlug"`
		State        flow2.State `bson:"state"`
		OccurredErrs []string    `bson:"occurredErrs"`
	}

	scenarioSlugDocument struct {
		Story    string `bson:"story"`
		Scenario string `bson:"scenario"`
	}
)

func newFlowDocument(flow *flow2.Flow) flowDocument {
	return flowDocument{
		ID:            flow.ID(),
		PerformanceID: flow.PerformanceID(),
		OverallState:  flow.OverallState(),
		Statuses:      newStatusDocuments(flow.Statuses()),
	}
}

func newStatusDocuments(statuses []*flow2.Status) []statusDocument {
	documents := make([]statusDocument, 0, len(statuses))
	for _, s := range statuses {
		documents = append(documents, newStatusDocument(s))
	}

	return documents
}

func newStatusDocument(status *flow2.Status) statusDocument {
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

func newThesisStatusDocuments(statuses []*flow2.ThesisStatus) []thesisStatusDocument {
	documents := make([]thesisStatusDocument, 0, len(statuses))
	for _, s := range statuses {
		documents = append(documents, newThesisStatusDocument(s))
	}

	return documents
}

func newThesisStatusDocument(status *flow2.ThesisStatus) thesisStatusDocument {
	return thesisStatusDocument{
		ThesisSlug:   status.ThesisSlug(),
		State:        status.State(),
		OccurredErrs: status.OccurredErrs(),
	}
}

func newFlow(d flowDocument) *flow2.Flow {
	return flow2.Unmarshal(flow2.Params{
		ID:            d.ID,
		PerformanceID: d.PerformanceID,
		OverallState:  d.OverallState,
		Statuses:      newStatuses(d.Statuses),
	})
}

func newStatuses(ds []statusDocument) []*flow2.Status {
	statuses := make([]*flow2.Status, 0, len(ds))
	for _, d := range ds {
		statuses = append(statuses, newStatus(d))
	}

	return statuses
}

func newStatus(d statusDocument) *flow2.Status {
	return flow2.NewStatus(
		newScenarioSlug(d.Slug),
		d.State,
		newThesisStatuses(d.ThesisStatuses)...,
	)
}

func newScenarioSlug(d scenarioSlugDocument) specification.Slug {
	return specification.NewScenarioSlug(d.Story, d.Scenario)
}

func newThesisStatuses(ds []thesisStatusDocument) []*flow2.ThesisStatus {
	statuses := make([]*flow2.ThesisStatus, 0, len(ds))
	for _, d := range ds {
		statuses = append(statuses, flow2.NewThesisStatus(
			d.ThesisSlug,
			d.State,
			d.OccurredErrs...,
		))
	}

	return statuses
}
