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

type ProductRepositoryInterface interface {
	Create(ctx context.Context, product *Product) error
	Find(ctx context.Context, ean, ncm string) (*[]Product, error)
	FindById(ctx context.Context, id string) (*Product, error)
}

type ProductRepository struct {
	Collection *mongo.Collection
	Logger     *slog.Logger
}

func NewRepository(client *mongo.Client, cfg *config.Configuration, logger *slog.Logger) *ProductRepository {
	collection := client.Database(cfg.DB).Collection(cfg.ProductCollection)
	return &ProductRepository{Collection: collection, Logger: logger}
}

func (pr *ProductRepository) Create(ctx context.Context, product *Product) error {
	_, err := pr.Collection.InsertOne(ctx, product)

	if err != nil {
		return err
	}

	return nil
}

func (pr *ProductRepository) Find(ctx context.Context, ean, ncm string) (*[]Product, error) {
	filter := bson.M{}

	if ean != "" {
		ean, _ := strconv.Atoi(ean)
		filter["ean"] = bson.M{"$eq": ean}
	}

	if ncm != "" {
		filter["ncm"] = bson.M{"$eq": ncm}
	}
	cursor, err := pr.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	products := []Product{}

	for cursor.Next(context.TODO()) {
		var p Product
		if err = cursor.Decode(&p); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return &products, nil
}
func (pr *ProductRepository) FindById(ctx context.Context, id string) (*Product, error) {
	filter := bson.M{
		"id": id,
	}
	cursor := pr.Collection.FindOne(ctx, filter)
	var product Product
	if err := cursor.Decode(&product); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &product, nil
}
