package mongodb

import (
	"github.com/harpyd/thestis/internal/core/entity/flow"
	"github.com/harpyd/thestis/internal/core/entity/specification"
)

type (
	flowDocument struct {
		ID           string           `bson:"_id"`
		PipelineID   string           `bson:"pipelineId"`
		OverallState flow.State       `bson:"overallState"`
		Statuses     []statusDocument `bson:"statuses"`
	}

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
		ID:           flow.ID(),
		PipelineID:   flow.PipelineID(),
		OverallState: flow.OverallState(),
		Statuses:     newStatusDocuments(flow.Statuses()),
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

func newFlow(d flowDocument) *flow.Flow {
	return flow.FromStatuses(d.ID, d.PipelineID, newStatuses(d.Statuses)...)
}

func newStatuses(ds []statusDocument) []*flow.Status {
	statuses := make([]*flow.Status, 0, len(ds))
	for _, d := range ds {
		statuses = append(statuses, newStatus(d))
	}

	return statuses
}

func newStatus(d statusDocument) *flow.Status {
	return flow.NewStatus(
		newScenarioSlug(d.Slug),
		d.State,
		newThesisStatuses(d.ThesisStatuses)...,
	)
}

func newScenarioSlug(d scenarioSlugDocument) specification.Slug {
	return specification.NewScenarioSlug(d.Story, d.Scenario)
}

func newThesisStatuses(ds []thesisStatusDocument) []*flow.ThesisStatus {
	statuses := make([]*flow.ThesisStatus, 0, len(ds))
	for _, d := range ds {
		statuses = append(statuses, flow.NewThesisStatus(
			d.ThesisSlug,
			d.State,
			d.OccurredErrs...,
		))
	}

	return statuses
}
