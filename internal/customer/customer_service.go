package customer

import (
	"context"
)

type CustomerService struct {
	Repository CustomerRepository
}

func NewService(cr *CustomerRepository) *CustomerService {
	return &CustomerService{
		Repository: *cr,
	}
}
func (service *CustomerService) Find(ctx context.Context, query FindQueryParams) (*[]Customer, error) {
	cursor, err := service.Repository.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	Customers := []Customer{}

	for cursor.Next(context.TODO()) {
		var p Customer
		if err = cursor.Decode(&p); err != nil {
			return nil, err
		}
		Customers = append(Customers, p)
	}
	return &Customers, nil
}
