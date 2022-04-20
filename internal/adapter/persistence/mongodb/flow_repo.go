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

type FlowRepository struct {
	flows *mongo.Collection
}

const flows = "flows"

func NewFlowRepository(db *mongo.Database) *FlowRepository {
	return &FlowRepository{
		flows: db.Collection(flows),
	}
}

func (r *FlowRepository) GetFlow(ctx context.Context, flowID string) (*flow.Flow, error) {
	document, err := r.getFlowDocument(ctx, bson.M{"_id": flowID})
	if err != nil {
		return nil, err
	}

	return newFlow(document), err
}

func (r *FlowRepository) getFlowDocument(ctx context.Context, filter bson.M) (flowDocument, error) {
	var document flowDocument
	if err := r.flows.FindOne(ctx, filter).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments); err != nil {
			return flowDocument{}, app.ErrFlowNotFound
		}

		return flowDocument{}, app.WrapWithDatabaseError(err)
	}

	return document, nil
}

func (r *FlowRepository) UpsertFlow(ctx context.Context, flow *flow.Flow) error {
	document := newFlowDocument(flow)

	opt := options.Replace().SetUpsert(true)
	_, err := r.flows.ReplaceOne(ctx, bson.M{"_id": flow.ID()}, document, opt)

	return app.WrapWithDatabaseError(err)
}
