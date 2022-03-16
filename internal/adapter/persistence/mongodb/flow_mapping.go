package mongodb

import (
	"github.com/harpyd/thestis/internal/domain/flow"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type (
	flowDocument struct {
		ID            string          `bson:"_id"`
		PerformanceID string          `bson:"performanceId"`
		Statuses      statusDocuments `bson:"statuses"`
	}

	statusDocuments []statusDocument

	statusDocument struct {
		Slug         slugDocument `bson:"slug"`
		State        flow.State   `bson:"state"`
		OccurredErrs []string     `bson:"occurredErrs"`
	}

	slugDocument struct {
		Kind     specification.SlugKind
		Story    string
		Scenario string
		Thesis   string
	}
)

func marshalToFlowDocument(flow flow.Flow) flowDocument {
	return flowDocument{
		ID:            flow.ID(),
		PerformanceID: flow.PerformanceID(),
		Statuses:      marshalToStatusDocuments(flow.Statuses()),
	}
}

func marshalToStatusDocuments(statuses []flow.Status) []statusDocument {
	documents := make([]statusDocument, 0, len(statuses))
	for _, s := range statuses {
		documents = append(documents, marshalToStatusDocument(s))
	}

	return documents
}

func marshalToStatusDocument(status flow.Status) statusDocument {
	return statusDocument{
		Slug:         marshalToSlugDocument(status.Slug()),
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

func (ds statusDocuments) unmarshalToStatuses() []flow.Status {
	transitions := make([]flow.Status, 0, len(ds))
	for _, d := range ds {
		transitions = append(transitions, d.unmarshalToStatus())
	}

	return transitions
}

func (d statusDocument) unmarshalToStatus() flow.Status {
	var slug specification.Slug

	switch d.Slug.Kind {
	case specification.StorySlug:
		slug = specification.NewStorySlug(d.Slug.Story)
	case specification.ScenarioSlug:
		slug = specification.NewScenarioSlug(d.Slug.Story, d.Slug.Scenario)
	case specification.ThesisSlug:
		slug = specification.NewThesisSlug(d.Slug.Story, d.Slug.Scenario, d.Slug.Thesis)
	case specification.NoSlug:
	}

	return flow.NewStatus(slug, d.State, d.OccurredErrs...)
}
