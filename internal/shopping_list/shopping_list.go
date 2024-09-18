package shoppinglist

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mercadola/api/pkg/exceptions"
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

type ShoppingListDto struct {
	Name        string   `json:"name" validate:"required"`
	CustomerId  string   `json:"customer_id" validate:"required"`
	ProductsIds []string `json:"products_ids"`
}

func (dto *ShoppingListDto) Validate() error {
	return exceptions.ValidateException(validator.New(), dto)
}
