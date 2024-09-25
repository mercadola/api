package customer

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mercadola/api/pkg/exceptions"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticateInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (dto *AuthenticateInput) Validate() error {
	return exceptions.ValidateException(validator.New(), dto)
}

type AutenticateOutput struct {
	AccessToken string `json:"accessToken"`
}

type FindByEmailInput struct {
	Email string `json:"email" validate:"required,email"`
}

func (dto *FindByEmailInput) Validate() error {
	return exceptions.ValidateException(validator.New(), dto)
}

type CustomerInterface interface {
	New(customerDto *CustomerDto) (*Customer, error)
	validatePassword(password string) bool
}

type Customer struct {
	ID        string            `json:"id" bson:"id"`
	Name      string            `json:"name" bson:"name"`
	Email     string            `json:"email" bson:"email"`
	Password  string            `json:"-" bson:"password"`
	CPF       string            `json:"cpf" bson:"cpf"`
	Phone     string            `json:"phone" bson:"phone"`
	Cep       string            `json:"cep" bson:"cep"`
	Gender    GenderEnumeration `json:"gender" bson:"gender"`
	Birthday  time.Time         `json:"birthday" bson:"birthday"`
	Active    bool              `json:"active" bson:"active,default=true"`
	CreatedAt time.Time         `json:"create_at" bson:"created_at"`
	UpdatedAt time.Time         `json:"update_at" bson:"updated_at"`
}

func (c *Customer) New(customerDto *CustomerDto) (*Customer, error) {
	pwHash, err := bcrypt.GenerateFromPassword([]byte(customerDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	c.ID = uuid.New().String()
	c.Name = customerDto.Name
	c.Email = customerDto.Email
	c.Password = string(pwHash)
	c.CPF = customerDto.CPF
	c.Phone = "+55" + customerDto.Phone
	c.Gender = customerDto.Gender
	c.Cep = customerDto.Cep
	c.Birthday = customerDto.Birthday
	c.Active = true
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	return c, nil
}

func (c *Customer) validatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password))

	return err == nil
}

type FindQueryParams struct {
	Email string
	CPF   string
}

func (params FindQueryParams) Validate() error {
	if params.Email == "" && params.CPF == "" {
		return exceptions.NewAppException(http.StatusBadRequest, "Email or CPF must be informed", nil)
	}
	return nil
}

type GenderEnumeration string

const (
	Male      GenderEnumeration = "Male"
	Female    GenderEnumeration = "Female"
	Undefined GenderEnumeration = "Undefined"
)

type CustomerDto struct {
	Name     string            `json:"name" validate:"required"`
	Email    string            `json:"email" validate:"required,email"`
	Password string            `json:"password" validate:"required,min=8,max=20"`
	CPF      string            `json:"cpf"`
	Phone    string            `json:"phone" validate:"required,len=11"`
	Cep      string            `json:"cep,omitempty" validate:"len=8"`
	Gender   GenderEnumeration `json:"gender,omitempty" validate:"oneof=Male Female Undefined"`
	Birthday time.Time         `json:"birthday,omitempty"`
}

func (dto *CustomerDto) Validate() error {
	return exceptions.ValidateException(validator.New(), dto)
}
