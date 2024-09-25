package product

import (
	"context"
	"log/slog"
)

type ProductService struct {
	Repository ProductRepository
	Logger     *slog.Logger
}

func NewService(pr *ProductRepository, logger *slog.Logger) *ProductService {
	return &ProductService{
		Repository: *pr,
	}
}
func (service *ProductService) Find(ctx context.Context, query FindProductQueryParams) (*[]Product, error) {
	cursor, err := service.Repository.Find(ctx, query)
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
