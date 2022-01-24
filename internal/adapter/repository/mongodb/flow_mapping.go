package mongodb

import "github.com/harpyd/thestis/internal/domain/performance"

type flowDocument struct {
}

func marshalToFlowDocument(flow performance.Flow) flowDocument {
	return flowDocument{}
}

func (d flowDocument) unmarshalToFlow() performance.Flow {
	return performance.Flow{}
}
