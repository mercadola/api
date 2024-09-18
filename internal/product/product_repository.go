package product

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/mercadola/api/internal/infrastruture/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FindProductQueryParams struct {
	Ean string
	Ncm string
}

type ProductRepository struct {
	Collection *mongo.Collection
	Logger     *slog.Logger
}

func NewRepository(client *mongo.Client, cfg *config.Configuration, logger *slog.Logger) *ProductRepository {
	collection := client.Database(cfg.DB).Collection(cfg.ProductCollection)
	return &ProductRepository{Collection: collection, Logger: logger}
}

func (pr *ProductRepository) Find(ctx context.Context, query FindProductQueryParams) (*mongo.Cursor, error) {
	filter := bson.M{}

	if query.Ean != "" {
		ean, _ := strconv.Atoi(query.Ean)
		filter["gtin"] = bson.M{"$eq": ean}
	}

	if query.Ncm != "" {
		filter["ncm.code"] = bson.M{"$eq": query.Ncm}
	}
	cursor, err := pr.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}
