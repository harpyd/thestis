package mongodb

import "github.com/harpyd/thestis/internal/domain/performance"

type (
	flowDocument struct {
		ID            string               `bson:"_id"`
		PerformanceID string               `bson:"performanceID"`
		State         performance.State    `bson:"state"`
		Transitions   []transitionDocument `bson:"transitions"`
	}

	transitionDocument struct {
		From         string            `bson:"from"`
		To           string            `bson:"to"`
		State        performance.State `bson:"state"`
		OccurredErrs []string          `bson:"occurredErrs"`
	}
)

func marshalToFlowDocument(flow performance.Flow) flowDocument {
	return flowDocument{
		ID:            flow.ID(),
		PerformanceID: flow.PerformanceID(),
		State:         flow.State(),
		Transitions:   marshalToTransitionDocuments(flow.Transitions()),
	}
}

func marshalToTransitionDocuments(transitions []performance.Transition) []transitionDocument {
	documents := make([]transitionDocument, 0, len(transitions))
	for _, t := range transitions {
		documents = append(documents, marshalToTransitionDocument(t))
	}

	return documents
}

func marshalToTransitionDocument(transition performance.Transition) transitionDocument {
	return transitionDocument{
		From:         transition.From(),
		To:           transition.To(),
		State:        transition.State(),
		OccurredErrs: transition.OccurredErrs(),
	}
}

func (d flowDocument) unmarshalToFlow() performance.Flow {
	return performance.UnmarshalFlow(performance.FlowParams{
		ID:            d.ID,
		PerformanceID: d.PerformanceID,
		State:         d.State,
		Transitions:   unmarshalToTransitions(d.Transitions),
	})
}

func unmarshalToTransitions(documents []transitionDocument) []performance.Transition {
	transitions := make([]performance.Transition, 0, len(documents))
	for _, d := range documents {
		transitions = append(transitions, d.unmarshalToTransition())
	}

	return transitions
}

func (d transitionDocument) unmarshalToTransition() performance.Transition {
	return performance.NewTransition(d.State, d.From, d.To, d.OccurredErrs...)
}
