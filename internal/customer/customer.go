package customer

import (
	"net/http"
	"time"

	"github.com/mercadola/api/internal/shared/utils/exceptions"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	CPF      string             `json:"cpf" bson:"cpf"`
	Phone    string             `json:"phone" bson:"phone"`
	Cep      string             `json:"cep" bson:"cep"`
	CreateAt time.Time          `json:"create_at" bson:"create_at"`
	UpdateAt time.Time          `json:"update_at" bson:"update_at"`
}

type findQueryParams struct {
	Name  string
	Email string
	CPF   string
}

func (params findQueryParams) Validate() error {
	if params.Name == "" && params.Email == "" && params.CPF == "" {
		return exceptions.NewAppException(http.StatusBadRequest, "Bad Request", "Name, Email or CPF must be informed", nil)
	}
	return nil
}

type CustomerDto struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
	CPF      string `json:"cpf" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Cep      string `json:"cep,omitempty"`
}

func (dto *CustomerDto) Validate() error {
	return dto.Validate()
}
