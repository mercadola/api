package shoppinglist

import (
	"context"
)

type ShoppingListService struct {
	Repository ShoppingListRepository
}

func NewService(slr *ShoppingListRepository) *ShoppingListService {
	return &ShoppingListService{
		Repository: *slr,
	}
}
func (service *ShoppingListService) FindByCustomerId(ctx context.Context, customer_id string) (*[]ShoppingList, error) {
	cursor, err := service.Repository.FindByCustomerId(ctx, customer_id)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	shoppingList := []ShoppingList{}

	for cursor.Next(context.TODO()) {
		var sl ShoppingList
		if err = cursor.Decode(&sl); err != nil {
			return nil, err
		}
		shoppingList = append(shoppingList, sl)
	}
	return &shoppingList, nil
}
