package shoppinglist

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ShoppingListMock struct {
	mock.Mock
}

func (m *ShoppingListMock) New(name, customer_id string, productsIds []string) *ShoppingList {
	args := m.Called(name, customer_id, productsIds)
	return args.Get(0).(*ShoppingList)
}

type ShoppingListRepositoryMock struct {
	mock.Mock
}

func (m *ShoppingListRepositoryMock) Create(ctx context.Context, shoppingList *ShoppingList) error {
	args := m.Called(ctx, shoppingList)
	return args.Error(0)
}

func (m *ShoppingListRepositoryMock) UpdateName(ctx context.Context, name, customer_id, shopping_list_id string) (*UpdateResult, error) {
	args := m.Called(ctx, name, customer_id, shopping_list_id)
	return args.Get(0).(*UpdateResult), args.Error(1)
}

func (m *ShoppingListRepositoryMock) UpdateProducts(ctx context.Context, customer_id, shopping_list_id string, products []string) (*UpdateResult, error) {
	args := m.Called(ctx, customer_id, shopping_list_id, products)
	return args.Get(0).(*UpdateResult), args.Error(1)
}

func (m *ShoppingListRepositoryMock) FindByCustomerId(ctx context.Context, customer_id string) (*[]ShoppingList, error) {
	args := m.Called(ctx, customer_id)
	return args.Get(0).(*[]ShoppingList), args.Error(1)
}

func (m *ShoppingListRepositoryMock) Delete(ctx context.Context, customer_id, shopping_list_id string) (*DeleteResult, error) {
	args := m.Called(ctx, customer_id, shopping_list_id)
	return args.Get(0).(*DeleteResult), args.Error(1)
}
