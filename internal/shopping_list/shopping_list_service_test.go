package shoppinglist

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	t.Run("should create a shopping list", func(t *testing.T) {
		shoppingListDto := &ShoppingListCreateDto{
			Name:        "Supermercado",
			ProductsIds: []string{},
		}
		customer_id := "any_customer_id"
		mockShoppingList := &ShoppingListMock{}
		expectedShoppingList := &ShoppingList{
			ID:          "mocked-uuid",
			Name:        shoppingListDto.Name,
			CustomerId:  customer_id,
			ProductsIds: shoppingListDto.ProductsIds,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		mockShoppingList.On("New", shoppingListDto.Name, customer_id, shoppingListDto.ProductsIds).Return(expectedShoppingList).Once()

		mockRepository := &ShoppingListRepositoryMock{}
		mockRepository.On("Create", context.Background(), expectedShoppingList).Return(nil)
		service := NewService(mockRepository, mockShoppingList)
		result, err := service.Create(context.Background(), shoppingListDto, customer_id)

		assert.Nil(t, err)
		assert.Equal(t, expectedShoppingList, result)
	})
	t.Run("should return error when create a shopping list", func(t *testing.T) {
		shoppingListDto := &ShoppingListCreateDto{
			Name:        "Supermercado",
			ProductsIds: []string{},
		}
		customer_id := "any_customer_id"
		mockShoppingList := &ShoppingListMock{}
		expectedShoppingList := &ShoppingList{
			ID:          "mocked-uuid",
			Name:        shoppingListDto.Name,
			CustomerId:  customer_id,
			ProductsIds: shoppingListDto.ProductsIds,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		mockShoppingList.On("New", shoppingListDto.Name, customer_id, shoppingListDto.ProductsIds).Return(expectedShoppingList).Once()

		mockRepository := &ShoppingListRepositoryMock{}
		mockRepository.On("Create", context.Background(), expectedShoppingList).Return(errors.New("any error"))

		service := NewService(mockRepository, mockShoppingList)
		result, err := service.Create(context.Background(), shoppingListDto, customer_id)

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
	})
}

func TestUpdateName(t *testing.T) {
	t.Run("should update a shopping list name", func(t *testing.T) {
		mockShoppingList := &ShoppingListMock{}
		mockRepository := &ShoppingListRepositoryMock{}
		mockRepository.On("UpdateName", context.Background(), "Novo nome", "any_customer_id", "any_shopping_list_id").Return(&UpdateResult{ModifiedCount: 1}, nil)
		service := NewService(mockRepository, mockShoppingList)
		err := service.UpdateName(context.Background(), "Novo nome", "any_customer_id", "any_shopping_list_id")

		assert.Nil(t, err)
	})
	t.Run("should return error when could not update a shopping list name", func(t *testing.T) {
		mockShoppingList := &ShoppingListMock{}
		mockRepository := &ShoppingListRepositoryMock{}
		mockRepository.On("UpdateName", context.Background(), "Novo nome", "any_customer_id", "any_shopping_list_id").Return(&UpdateResult{}, errors.New("any error"))
		service := NewService(mockRepository, mockShoppingList)
		err := service.UpdateName(context.Background(), "Novo nome", "any_customer_id", "any_shopping_list_id")

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
	})
	t.Run("should return error when could not find a shopping list to update", func(t *testing.T) {
		mockShoppingList := &ShoppingListMock{}
		mockRepository := &ShoppingListRepositoryMock{}
		mockRepository.On("UpdateName", context.Background(), "Novo nome", "any_customer_id", "any_shopping_list_id").Return(&UpdateResult{ModifiedCount: 0}, nil)
		service := NewService(mockRepository, mockShoppingList)
		err := service.UpdateName(context.Background(), "Novo nome", "any_customer_id", "any_shopping_list_id")

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
	})
}

func TestUpdateProducts(t *testing.T) {
	t.Run("should update a shopping list name", func(t *testing.T) {
		mockShoppingList := &ShoppingListMock{}
		mockRepository := &ShoppingListRepositoryMock{}
		products := []string{"any_product_id"}
		mockRepository.On("UpdateProducts", context.Background(), "any_customer_id", "any_shopping_list_id", products).Return(&UpdateResult{ModifiedCount: 1}, nil)
		service := NewService(mockRepository, mockShoppingList)
		err := service.UpdateProducts(context.Background(), "any_customer_id", "any_shopping_list_id", products)

		assert.Nil(t, err)
	})
	t.Run("should return error when could not update a shopping list name", func(t *testing.T) {
		mockShoppingList := &ShoppingListMock{}
		mockRepository := &ShoppingListRepositoryMock{}
		products := []string{"any_product_id"}
		mockRepository.On("UpdateProducts", context.Background(), "any_customer_id", "any_shopping_list_id", products).Return(&UpdateResult{}, errors.New("any error"))
		service := NewService(mockRepository, mockShoppingList)
		err := service.UpdateProducts(context.Background(), "any_customer_id", "any_shopping_list_id", products)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
	})
	t.Run("should return error when could not find a shopping list to update", func(t *testing.T) {
		mockShoppingList := &ShoppingListMock{}
		mockRepository := &ShoppingListRepositoryMock{}
		products := []string{"any_product_id"}
		mockRepository.On("UpdateProducts", context.Background(), "any_customer_id", "any_shopping_list_id", products).Return(&UpdateResult{ModifiedCount: 0}, nil)
		service := NewService(mockRepository, mockShoppingList)
		err := service.UpdateProducts(context.Background(), "any_customer_id", "any_shopping_list_id", products)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
	})
}

func TestFindByCustomerId(t *testing.T) {
	t.Run("should find shopping lists by customer id", func(t *testing.T) {
		shoppingList := &ShoppingList{
			ID:         "any_id",
			CustomerId: "any_customer_id",
			Name:       "Supermercado",
		}
		result := &[]ShoppingList{*shoppingList}
		mockRepository := &ShoppingListRepositoryMock{}
		mockShoppingList := &ShoppingListMock{}
		mockRepository.On("FindByCustomerId", context.Background(), "any_customer_id").Return(result, nil)
		service := NewService(mockRepository, mockShoppingList)
		result, err := service.FindByCustomerId(context.Background(), "any_customer_id")

		assert.Nil(t, err)
		assert.Equal(t, 1, len(*result))
	})
	t.Run("should return error when find shopping lists by customer id", func(t *testing.T) {
		result := &[]ShoppingList{}
		mockRepository := &ShoppingListRepositoryMock{}
		mockRepository.On("FindByCustomerId", context.Background(), "any_customer_id").Return(result, errors.New("any error"))
		mockShoppingList := &ShoppingListMock{}
		service := NewService(mockRepository, mockShoppingList)
		result, err := service.FindByCustomerId(context.Background(), "any_customer_id")

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
	})
}

func TestDelete(t *testing.T) {
	t.Run("should delete a shopping list", func(t *testing.T) {
		mockShoppingList := &ShoppingList{}
		mockRepository := &ShoppingListRepositoryMock{}
		mockRepository.On("Delete", context.Background(), "any_customer_id", "any_shopping_list_id").Return(&DeleteResult{DeletedCount: 1}, nil)
		service := NewService(mockRepository, mockShoppingList)
		err := service.Delete(context.Background(), "any_customer_id", "any_shopping_list_id")

		assert.Nil(t, err)
	})
	t.Run("should return error when could not delete a shopping list", func(t *testing.T) {
		mockRepository := &ShoppingListRepositoryMock{}
		mockShoppingList := &ShoppingList{}
		mockRepository.On("Delete", context.Background(), "any_customer_id", "any_shopping_list_id").Return(&DeleteResult{}, errors.New("any error"))
		service := NewService(mockRepository, mockShoppingList)
		err := service.Delete(context.Background(), "any_customer_id", "any_shopping_list_id")

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
	})
	t.Run("should return error when could not find a shopping list to delete", func(t *testing.T) {
		mockRepository := &ShoppingListRepositoryMock{}
		mockShoppingList := &ShoppingList{}
		mockRepository.On("Delete", context.Background(), "any_customer_id", "any_shopping_list_id").Return(&DeleteResult{DeletedCount: 0}, nil)
		service := NewService(mockRepository, mockShoppingList)
		err := service.Delete(context.Background(), "any_customer_id", "any_shopping_list_id")

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
	})
}
