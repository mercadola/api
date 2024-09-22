package shoppinglist

import (
	"context"
	"log/slog"
	"time"

	"github.com/mercadola/api/internal/infrastruture/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UpdateNameResult struct {
	ModifiedCount int64 `json:"modified_count"`
}

type DeleteResult struct {
	DeletedCount int64 `json:"deleted_count"`
}

type ShoppingListRepositoryInterface interface {
	Create(ctx context.Context, shoppingList *ShoppingList) error
	UpdateName(ctx context.Context, name, customer_id, shopping_list_id string) (*UpdateNameResult, error)
	FindByCustomerId(ctx context.Context, customer_id string) (*[]ShoppingList, error)
	Delete(ctx context.Context, customer_id, shopping_list_id string) (*DeleteResult, error)
}

type ShoppingListRepository struct {
	Collection *mongo.Collection
	Logger     *slog.Logger
}

func NewRepository(client database.MongoClientInterface, logger *slog.Logger, database, shoppingListCollection string) *ShoppingListRepository {
	collection := client.Database(database).Collection(shoppingListCollection)
	return &ShoppingListRepository{Collection: collection, Logger: logger}
}

func (slr *ShoppingListRepository) Create(ctx context.Context, shoppingList *ShoppingList) error {
	_, err := slr.Collection.InsertOne(ctx, shoppingList)

	if err != nil {
		return err
	}

	return nil
}
func (slr *ShoppingListRepository) UpdateName(ctx context.Context, name, customer_id, shopping_list_id string) (*UpdateNameResult, error) {
	filter := bson.M{
		"customer_id": customer_id,
		"id":          shopping_list_id,
	}
	update := bson.M{
		"$set": bson.M{
			"name":       name,
			"updated_at": time.Now(),
		},
	}

	result, err := slr.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return &UpdateNameResult{ModifiedCount: result.ModifiedCount}, nil
}

func (slr *ShoppingListRepository) FindByCustomerId(ctx context.Context, customer_id string) (*[]ShoppingList, error) {
	filter := bson.M{
		"customer_id": customer_id,
	}
	cursor, err := slr.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	shoppingList := []ShoppingList{}

	for cursor.Next(ctx) {
		var sl ShoppingList
		if err = cursor.Decode(&sl); err != nil {
			return nil, err
		}
		shoppingList = append(shoppingList, sl)
	}
	return &shoppingList, nil
}

func (slr *ShoppingListRepository) Delete(ctx context.Context, customer_id, shopping_list_id string) (*DeleteResult, error) {
	filter := bson.M{
		"customer_id": customer_id,
		"id":          shopping_list_id,
	}
	result, err := slr.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &DeleteResult{DeletedCount: result.DeletedCount}, nil

}
