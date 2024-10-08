package customer

import (
	"context"
	"log/slog"

	"github.com/mercadola/api/internal/infrastruture/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerRepository struct {
	Collection *mongo.Collection
	Logger     *slog.Logger
}

func NewCustomerRepository(client *mongo.Client, cfg *config.Configuration, logger *slog.Logger) *CustomerRepository {
	collection := client.Database(cfg.DB).Collection(cfg.CustomerCollection)
	return &CustomerRepository{Collection: collection, Logger: logger}
}

func (cr *CustomerRepository) Create(ctx context.Context, customer Customer) error {
	_, err := cr.Collection.InsertOne(ctx, customer)

	if err != nil {
		return err
	}

	return nil
}

func (cr *CustomerRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := cr.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (cr *CustomerRepository) Find(ctx context.Context, query findQueryParams) (*mongo.Cursor, error) {
	params := bson.A{}

	if query.Email != "" {
		params = append(params, bson.M{"email": query.Email})
	}

	if query.CPF != "" {
		params = append(params, bson.M{"cpf": query.CPF})
	}

	filter := bson.M{"$or": params}

	cursor, err := cr.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func (cr *CustomerRepository) FindByEmail(ctx context.Context, email string) *mongo.SingleResult {
	filter := bson.M{"email": email}
	return cr.Collection.FindOne(ctx, filter)
}

func (cr *CustomerRepository) FindById(ctx context.Context, id primitive.ObjectID) *mongo.SingleResult {
	filter := bson.M{"_id": id}
	return cr.Collection.FindOne(ctx, filter)
}

func (cr *CustomerRepository) Update(ctx context.Context, customer Customer) error {
	filter := bson.M{"_id": customer.ID}
	update := bson.M{"$set": customer}
	_, err := cr.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
