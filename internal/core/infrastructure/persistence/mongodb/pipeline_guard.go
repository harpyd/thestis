package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/pipeline"
)

type PipelineGuard struct {
	pipelines *mongo.Collection
}

func NewPipelineGuard(db *mongo.Database) *PipelineGuard {
	return &PipelineGuard{
		pipelines: db.Collection(pipelineCollection),
	}
}

func (g *PipelineGuard) AcquirePipeline(ctx context.Context, pipeID string) error {
	var document pipelineDocument

	var (
		filter = bson.M{"_id": pipeID}
		update = bson.M{"$set": bson.M{"started": true}}
		opt    = options.FindOneAndUpdate().SetProjection(bson.M{"_id": 0, "started": 1})
	)

	if err := g.pipelines.FindOneAndUpdate(ctx, filter, update, opt).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return service.ErrPipelineNotFound
		}

		return service.WrapWithDatabaseError(err)
	}

	if document.Started {
		return pipeline.ErrAlreadyStarted
	}

	return nil
}

func (g *PipelineGuard) ReleasePipeline(ctx context.Context, pipeID string) error {
	update := bson.M{"$set": bson.M{"started": false}}

	_, err := g.pipelines.UpdateByID(ctx, pipeID, update)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return service.ErrPipelineNotFound
	}

	return service.WrapWithDatabaseError(err)
}
