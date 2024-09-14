package product

import (
	"context"
)

type ProductService struct {
	Repository ProductRepository
}

func NewService(pr *ProductRepository) *ProductService {
	return &ProductService{
		Repository: *pr,
	}
}
func (service *ProductService) Find(ctx context.Context, query FindQueryParams) (*[]Product, error) {
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
