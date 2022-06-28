package mongodb

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/pipeline"
)

type PipelineRepository struct {
	pipelines *mongo.Collection
}

const pipelineCollection = "pipelines"

func NewPipelineRepository(db *mongo.Database) *PipelineRepository {
	r := &PipelineRepository{
		pipelines: db.Collection(pipelineCollection),
	}

	return r
}

func (r *PipelineRepository) GetPipeline(
	ctx context.Context,
	pipeID string,
	specGetter service.SpecificationGetter,
	registrars ...pipeline.ExecutorRegistrar,
) (*pipeline.Pipeline, error) {
	document, err := r.getPipelineDocument(ctx, bson.M{"_id": pipeID})
	if err != nil {
		return nil, err
	}

	spec, err := specGetter.GetSpecification(ctx, document.SpecificationID)
	if err != nil {
		return nil, err
	}

	return newPipeline(document, spec, registrars), nil
}

func (r *PipelineRepository) getPipelineDocument(
	ctx context.Context,
	filter bson.M,
) (pipelineDocument, error) {
	var document pipelineDocument
	if err := r.pipelines.FindOne(ctx, filter).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return pipelineDocument{}, service.ErrPipelineNotFound
		}

		return pipelineDocument{}, service.WrapWithDatabaseError(err)
	}

	return document, nil
}

func (r *PipelineRepository) AddPipeline(ctx context.Context, pipe *pipeline.Pipeline) error {
	_, err := r.pipelines.InsertOne(ctx, newPipelineDocument(pipe))

	return service.WrapWithDatabaseError(err)
}
