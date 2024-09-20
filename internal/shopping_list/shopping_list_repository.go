package shoppinglist

import (
	"context"
	"log/slog"
	"time"

	"github.com/mercadola/api/internal/infrastruture/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (slr *ShoppingListRepository) Create(ctx context.Context, shoppingList *ShoppingList) error {
	_, err := slr.Collection.InsertOne(ctx, shoppingList)

	if err != nil {
		return err
	}

	return nil
}
func (slr *ShoppingListRepository) UpdateName(ctx context.Context, customer_id, shopping_list_id, name string) error {
	objID, err := primitive.ObjectIDFromHex(shopping_list_id)
	if err != nil {
		return err
	}

	filter := bson.M{
		"customer_id": customer_id,
		"_id":         objID,
	}
	update := bson.M{
		"$set": bson.M{
			"name":       name,
			"updated_at": time.Now(),
		},
	}

	_, err = slr.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
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
func (slr *ShoppingListRepository) FindById(ctx context.Context, customer_id, shopping_list_id string) (*mongo.SingleResult, error) {
	objID, err := primitive.ObjectIDFromHex(shopping_list_id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id":         objID,
		"customer_id": customer_id,
	}
	cursor := slr.Collection.FindOne(ctx, filter)
	return cursor, nil
}

func (slr *ShoppingListRepository) Delete(ctx context.Context, customer_id, shopping_list_id string) error {
	objID, err := primitive.ObjectIDFromHex(shopping_list_id)
	if err != nil {
		return err
	}

	filter := bson.M{
		"customer_id": customer_id,
		"_id":         objID,
	}
	_, err = slr.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil

}
