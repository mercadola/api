package customer

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/mercadola/api/internal/infrastruture/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FindQueryParams struct {
	Ean string
	Ncm string
}

type CustomerRepository struct {
	Collection *mongo.Collection
	Logger     *slog.Logger
}

func NewCustomerRepository(client *mongo.Client, cfg *config.Configuration, logger *slog.Logger) *CustomerRepository {
	collection := client.Database(cfg.DB).Collection(cfg.CustomerCollection)
	return &CustomerRepository{Collection: collection, Logger: logger}
}

func (cr *CustomerRepository) Find(ctx context.Context, query FindQueryParams) (*mongo.Cursor, error) {
	filter := bson.M{}

	if query.Ean != "" {
		ean, _ := strconv.Atoi(query.Ean)
		filter["gtin"] = bson.M{"$eq": ean}
	}

	if query.Ncm != "" {
		filter["ncm.code"] = bson.M{"$eq": query.Ncm}
	}
	cursor, err := cr.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}
