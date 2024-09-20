package shoppinglist

import (
	"context"
	"net/http"

	"github.com/mercadola/api/pkg/exceptions"
)

type ShoppingListService struct {
	Repository ShoppingListRepository
}

func NewService(slr *ShoppingListRepository) *ShoppingListService {
	return &ShoppingListService{
		Repository: *slr,
	}
}

func (service *ShoppingListService) Create(ctx context.Context, shoppingListDto *ShoppingListCreateDto, customer_id string) (*ShoppingList, error) {
	shoppingList := NewShoppingList(shoppingListDto.Name, customer_id, shoppingListDto.ProductsIds)
	err := service.Repository.Create(ctx, shoppingList)
	if err != nil {
		return nil, exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}

	return shoppingList, nil
}
func (service *ShoppingListService) UpdateName(ctx context.Context, name, customer_id, shopping_list_id string) error {
	cursor, err := service.Repository.FindById(ctx, customer_id, shopping_list_id)
	if err != nil {
		return err
	}
	var sl ShoppingList
	if err = cursor.Decode(&sl); err != nil {
		return err
	}
	return nil

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

func (service *ShoppingListService) Delete(ctx context.Context, customer_id, shopping_list_id string) error {
	err := service.Repository.Delete(ctx, customer_id, shopping_list_id)
	if err != nil {
		return err
	}
	return nil

}
