package customer

import (
	"context"
	"log/slog"
	"time"

	"github.com/mercadola/api/internal/infrastruture/config"
	"github.com/mercadola/api/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerRepositoryInterface interface {
	Create(ctx context.Context, customer *Customer) error
	Delete(ctx context.Context, id string) (*utils.DeleteResult, error)
	Find(ctx context.Context, query *FindQueryParams) (*[]Customer, error)
	FindByEmail(ctx context.Context, findByEmailInput *FindByEmailInput) (*Customer, error)
	FindById(ctx context.Context, id string) (*Customer, error)
	InactiveCustomer(ctx context.Context, id string) (*utils.UpdateResult, error)
	PositivateCustomer(ctx context.Context, id string) (*utils.UpdateResult, error)
	Update(ctx context.Context, customer *Customer) (*utils.UpdateResult, error)
}

type CustomerRepository struct {
	Collection *mongo.Collection
	Logger     *slog.Logger
}

func NewCustomerRepository(client *mongo.Client, cfg *config.Configuration, logger *slog.Logger) *CustomerRepository {
	collection := client.Database(cfg.DB).Collection(cfg.CustomerCollection)
	return &CustomerRepository{Collection: collection, Logger: logger}
}

func (cr *CustomerRepository) Create(ctx context.Context, customer *Customer) error {
	_, err := cr.Collection.InsertOne(ctx, customer)

	if err != nil {
		return err
	}

	return nil
}

func (cr *CustomerRepository) Delete(ctx context.Context, id string) (*utils.DeleteResult, error) {
	filter := bson.M{"_id": id}
	result, err := cr.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &utils.DeleteResult{DeletedCount: result.DeletedCount}, nil
}

func (cr *CustomerRepository) Find(ctx context.Context, query *FindQueryParams) (*[]Customer, error) {
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
	defer cursor.Close(ctx)
	customers := []Customer{}

	for cursor.Next(ctx) {
		var c Customer
		if err = cursor.Decode(&c); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return &customers, nil
}

func (cr *CustomerRepository) FindByEmail(ctx context.Context, findByEmailInput *FindByEmailInput) (*Customer, error) {
	filter := bson.M{"email": findByEmailInput.Email}
	result := cr.Collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, result.Err()
	}
	var customer Customer

	err := result.Decode(&customer)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (cr *CustomerRepository) FindById(ctx context.Context, id string) (*Customer, error) {
	filter := bson.M{"_id": id}
	result := cr.Collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, result.Err()
	}
	var customer Customer

	err := result.Decode(&customer)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (cr *CustomerRepository) InactiveCustomer(ctx context.Context, id string) (*utils.UpdateResult, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"active": false, "updated_at": time.Now()}}
	result, err := cr.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return &utils.UpdateResult{ModifiedCount: result.ModifiedCount}, nil
}

func (cr *CustomerRepository) PositivateCustomer(ctx context.Context, id string) (*utils.UpdateResult, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"active": true, "updated_at": time.Now()}}
	result, err := cr.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return &utils.UpdateResult{ModifiedCount: result.ModifiedCount}, nil
}

func (cr *CustomerRepository) Update(ctx context.Context, customer *Customer) (*utils.UpdateResult, error) {
	filter := bson.M{"_id": customer.ID}
	update := bson.M{"$set": customer}
	result, err := cr.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return &utils.UpdateResult{ModifiedCount: result.ModifiedCount}, nil
}
