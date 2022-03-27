package mongodb

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/flow"
)

type FlowsRepository struct {
	flows *mongo.Collection
}

const flows = "flows"

func NewFlowsRepository(db *mongo.Database) *FlowsRepository {
	return &FlowsRepository{
		flows: db.Collection(flows),
	}
}

func (r *FlowsRepository) GetFlow(ctx context.Context, flowID string) (flow.Flow, error) {
	document, err := r.getFlowDocument(ctx, bson.M{"_id": flowID})
	if err != nil {
		return flow.Flow{}, err
	}

	return document.unmarshalToFlow(), err
}

func (r *FlowsRepository) getFlowDocument(ctx context.Context, filter bson.M) (flowDocument, error) {
	var document flowDocument
	if err := r.flows.FindOne(ctx, filter).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments); err != nil {
			return flowDocument{}, app.ErrFlowNotFound
		}

		return flowDocument{}, app.WrapWithDatabaseError(err)
	}

	return document, nil
}

func (r *FlowsRepository) UpsertFlow(ctx context.Context, flow flow.Flow) error {
	document := marshalToFlowDocument(flow)

	opt := options.Replace().SetUpsert(true)
	_, err := r.flows.ReplaceOne(ctx, bson.M{"_id": flow.ID()}, document, opt)

	return app.WrapWithDatabaseError(err)
}

func (r *FlowsRepository) RemoveAllFlows(ctx context.Context) error {
	_, err := r.flows.DeleteOne(ctx, bson.D{})

	return app.WrapWithDatabaseError(err)
}
