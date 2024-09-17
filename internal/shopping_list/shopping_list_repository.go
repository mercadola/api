package shoppinglist

import (
	"context"
	"log/slog"

	"github.com/mercadola/api/internal/infrastruture/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShoppingListRepository struct {
	Collection *mongo.Collection
	Logger     *slog.Logger
}

func NewRepository(client *mongo.Client, cfg *config.Configuration, logger *slog.Logger) *ShoppingListRepository {
	collection := client.Database(cfg.DB).Collection(cfg.ShoppingListCollection)
	return &ShoppingListRepository{Collection: collection, Logger: logger}
}

func (slr *ShoppingListRepository) FindByCustomerId(ctx context.Context, customer_id string) (*mongo.Cursor, error) {
	filter := bson.M{
		"customer_id": customer_id,
	}
	cursor, err := slr.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}
