package shoppinglist

import (
	"context"
	"net/http"

	"github.com/mercadola/api/pkg/exceptions"
)

type ShoppingListService struct {
	Repository   ShoppingListRepositoryInterface
	ShoppingList ShoppingListInterface
}

func NewService(slr ShoppingListRepositoryInterface, sl ShoppingListInterface) *ShoppingListService {
	return &ShoppingListService{
		Repository:   slr,
		ShoppingList: sl}
}

func (service *ShoppingListService) Create(ctx context.Context, shoppingListDto *ShoppingListCreateDto, customer_id string) (*ShoppingList, *exceptions.AppException) {
	shoppingList := service.ShoppingList.New(shoppingListDto.Name, customer_id, shoppingListDto.ProductsIds)
	err := service.Repository.Create(ctx, shoppingList)
	if err != nil {
		return nil, exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}

	return shoppingList, nil
}
func (service *ShoppingListService) UpdateName(ctx context.Context, name, customer_id, shopping_list_id string) *exceptions.AppException {
	result, err := service.Repository.UpdateName(ctx, name, customer_id, shopping_list_id)
	if err != nil {
		return exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}
	if result.ModifiedCount == 0 {
		return exceptions.NewAppException(http.StatusNotFound, "no documents in result", nil)
	}
	return nil

}

func (service *ShoppingListService) FindByCustomerId(ctx context.Context, customer_id string) (*[]ShoppingList, *exceptions.AppException) {
	shoppingList, err := service.Repository.FindByCustomerId(ctx, customer_id)
	if err != nil {
		return nil, exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}
	return shoppingList, nil
}

func (service *ShoppingListService) Delete(ctx context.Context, customer_id, shopping_list_id string) *exceptions.AppException {
	result, err := service.Repository.Delete(ctx, customer_id, shopping_list_id)
	if err != nil {
		return exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}
	if result.DeletedCount == 0 {
		return exceptions.NewAppException(http.StatusNotFound, "no documents in result", nil)
	}
	return nil

}
