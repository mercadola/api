package shoppinglist

import (
	"context"
	"net/http"
	"time"

	"github.com/mercadola/api/pkg/exceptions"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (service *ShoppingListService) Create(ctx context.Context, shoppingListDto *ShoppingListDto) (*ShoppingList, error) {
	if err := shoppingListDto.Validate(); err != nil {
		return nil, err
	}
	shoppingList := &ShoppingList{
		ID:          primitive.NewObjectID(),
		Name:        shoppingListDto.Name,
		CustomerId:  shoppingListDto.CustomerId,
		ProductsIds: shoppingListDto.ProductsIds,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := service.Repository.Create(ctx, shoppingList)
	if err != nil {
		return nil, exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}

	return shoppingList, nil
}
