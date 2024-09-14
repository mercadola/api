package product

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FindQueryParams struct {
	Ean string
	Ncm string
}

type ProductRepository struct {
	Collection *mongo.Collection
	Logger     *slog.Logger
}

func NewProductRepository(client *mongo.Client, logger *slog.Logger) *ProductRepository {
	collection := client.Database(os.Getenv("MERCADOLA_DATABASE")).Collection(os.Getenv("PRODUCT_COLLECTION"))
	return &ProductRepository{Collection: collection, Logger: logger}

}
func (pr *ProductRepository) Find(ctx context.Context, query FindQueryParams) (*mongo.Cursor, error) {
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
