package shoppinglist

import (
	"time"

	"github.com/google/uuid"
)

type ShoppingListInterface interface {
	New(name, customer_id string, productsIds []string) *ShoppingList
}

type ShoppingList struct {
	ID          string    `json:"id" bson:"id"`
	CustomerId  string    `json:"customer_id" bson:"customer_id"`
	Name        string    `json:"name" bson:"name"`
	ProductsIds []string  `json:"products_ids" bson:"products_ids"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

func (sl *ShoppingList) New(name, customer_id string, productsIds []string) *ShoppingList {
	sl.ID = uuid.New().String()
	sl.Name = name
	sl.CustomerId = customer_id
	sl.ProductsIds = productsIds
	sl.CreatedAt = time.Now()
	sl.UpdatedAt = time.Now()
	return sl
}

type ShoppingListCreateDto struct {
	Name        string   `json:"name" validate:"required"`
	ProductsIds []string `json:"products_ids"`
}

type ShoppingListUpdateDto struct {
	Name string `json:"name" validate:"required"`
}
