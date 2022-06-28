package mongodb

import (
	"github.com/harpyd/thestis/internal/core/entity/pipeline"
	"github.com/harpyd/thestis/internal/core/entity/specification"
)

type pipelineDocument struct {
	ID              string `bson:"_id"`
	OwnerID         string `bson:"ownerId"`
	SpecificationID string `bson:"specificationId"`
	Started         bool   `bson:"started"`
}

func newPipelineDocument(pipe *pipeline.Pipeline) pipelineDocument {
	return pipelineDocument{
		ID:              pipe.ID(),
		OwnerID:         pipe.OwnerID(),
		SpecificationID: pipe.SpecificationID(),
		Started:         pipe.Started(),
	}
}

func newPipeline(
	d pipelineDocument,
	spec *specification.Specification,
	registrars []pipeline.ExecutorRegistrar,
) *pipeline.Pipeline {
	return pipeline.Unmarshal(pipeline.Params{
		ID:            d.ID,
		Specification: spec,
		OwnerID:       d.OwnerID,
		Started:       d.Started,
	}, registrars...)
}
