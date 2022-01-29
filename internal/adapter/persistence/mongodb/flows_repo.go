package mongodb

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/performance"
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

func (r *FlowsRepository) GetFlow(ctx context.Context, flowID string) (performance.Flow, error) {
	document, err := r.getFlowDocument(ctx, flowID, "")
	if err != nil {
		return performance.Flow{}, err
	}

	return document.unmarshalToFlow(), err
}

func (r *FlowsRepository) getFlowDocument(ctx context.Context, flowID, userID string) (flowDocument, error) {
	filter := makeFlowFilter(flowID, userID)

	var document flowDocument
	if err := r.flows.FindOne(ctx, filter).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments); err != nil {
			return flowDocument{}, app.NewFlowNotFoundError(err)
		}

		return flowDocument{}, app.NewDatabaseError(err)
	}

	return document, nil
}

func makeFlowFilter(flowID string, userID string) bson.M {
	filter := bson.M{"_id": flowID}
	if userID != "" {
		filter["ownerId"] = userID
	}

	return filter
}

func (r *FlowsRepository) UpsertFlow(ctx context.Context, flow performance.Flow) error {
	document := marshalToFlowDocument(flow)

	opt := options.Replace().SetUpsert(true)
	_, err := r.flows.ReplaceOne(ctx, bson.M{"_id": flow.ID()}, document, opt)

	return app.NewDatabaseError(err)
}

func (r *FlowsRepository) RemoveAllFlows(ctx context.Context) error {
	_, err := r.flows.DeleteOne(ctx, bson.D{})

	return app.NewDatabaseError(err)
}
