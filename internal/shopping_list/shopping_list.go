package shoppinglist

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShoppingList struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	CustomerId  string             `json:"customer_id" bson:"customer_id"`
	Name        string             `json:"name" bson:"name"`
	ProductsIds []string           `json:"products_ids" bson:"products_ids"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

func NewShoppingList(name, customer_id string, productsIds []string) *ShoppingList {
	return &ShoppingList{
		ID:          primitive.NewObjectID(),
		Name:        name,
		CustomerId:  customer_id,
		ProductsIds: productsIds,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

type ShoppingListCreateDto struct {
	Name        string   `json:"name" validate:"required"`
	ProductsIds []string `json:"products_ids"`
}

type ShoppingListUpdateDto struct {
	Name string `json:"name" validate:"required"`
}
