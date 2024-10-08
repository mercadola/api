package customer

import (
	"context"
	"net/http"
	"time"

	"github.com/mercadola/api/pkg/exceptions"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type CustomerService struct {
	Repository CustomerRepository
}

func NewService(cr *CustomerRepository) *CustomerService {
	return &CustomerService{
		Repository: *cr,
	}
}

func (service *CustomerService) Authenticate(ctx context.Context, authenticateInput AuthenticateInput) (*Customer, error) {
	finddedCustomer, err := service.FindByEmail(ctx, FindByEmailInput{Email: authenticateInput.Email})

	if err != nil {
		return nil, exceptions.NewAppException(http.StatusUnauthorized, "Invalid credentials", nil)
	}

	if !finddedCustomer.validatePassword(authenticateInput.Password) {
		return nil, exceptions.NewAppException(http.StatusUnauthorized, "Invalid credentials", nil)
	}

	return finddedCustomer, nil
}

func (service *CustomerService) Create(ctx context.Context, customerDto *CustomerDto) (*Customer, error) {
	finddedCustomers, _ := service.Find(ctx, findQueryParams{Email: customerDto.Email, CPF: customerDto.CPF})

	if finddedCustomers != nil && len(*finddedCustomers) > 0 {
		return nil, exceptions.NewAppException(http.StatusConflict, "Customer already exists", nil)
	}

	pwHash, err := bcrypt.GenerateFromPassword([]byte(customerDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}

	customer := Customer{
		ID:        primitive.NewObjectID(),
		Name:      customerDto.Name,
		Email:     customerDto.Email,
		Password:  string(pwHash),
		CPF:       customerDto.CPF,
		Phone:     "+55" + customerDto.Phone,
		Cep:       customerDto.Cep,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = service.Repository.Create(ctx, customer)
	if err != nil {
		return nil, exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}

	// TODO DISPARO DE E-MAIL DE BOAS VINDAS

	return &customer, nil
}

func (service *CustomerService) Delete(ctx context.Context, id primitive.ObjectID) error {
	customer, err := service.FindById(ctx, id)
	if err != nil {
		return err
	}

	customer.CPF = ""
	customer.Active = false
	customer.UpdatedAt = time.Now()

	err = service.Repository.Update(ctx, *customer)

	return nil
}

func (service *CustomerService) Find(ctx context.Context, query findQueryParams) (*[]Customer, error) {
	cursor, err := service.Repository.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	customers := []Customer{}

	for cursor.Next(context.TODO()) {
		var p Customer
		if err = cursor.Decode(&p); err != nil {
			return nil, err
		}
		customers = append(customers, p)
	}
	return &customers, nil
}

func (service *CustomerService) FindByEmail(ctx context.Context, findByEmail FindByEmailInput) (*Customer, error) {
	result := service.Repository.FindByEmail(ctx, findByEmail.Email)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, exceptions.NewAppException(http.StatusNotFound, "Customer not found", nil)
		}
		return nil, exceptions.NewAppException(http.StatusInternalServerError, result.Err().Error(), nil)
	}

	var customer Customer
	result.Decode(&customer)

	return &customer, nil
}

func (service *CustomerService) FindById(ctx context.Context, id primitive.ObjectID) (*Customer, error) {
	result := service.Repository.FindById(ctx, id)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, exceptions.NewAppException(http.StatusNotFound, "Customer not found", nil)
		}
		return nil, exceptions.NewAppException(http.StatusInternalServerError, result.Err().Error(), nil)
	}

	var customer Customer
	result.Decode(&customer)

	return &customer, nil
}

func (service *CustomerService) PositivateCustomer(ctx context.Context, id primitive.ObjectID) error {
	customer, err := service.FindById(ctx, id)
	if err != nil {
		return err
	}

	customer.Active = true
	customer.UpdatedAt = time.Now()

	err = service.Repository.Update(ctx, *customer)
	if err != nil {
		return err
	}

	return nil
}

func (service *CustomerService) Update(ctx context.Context, id primitive.ObjectID, customerDto *CustomerDto) error {
	customer := Customer{
		ID:        id,
		Name:      customerDto.Name,
		Email:     customerDto.Email,
		Password:  customerDto.Password,
		CPF:       customerDto.CPF,
		Phone:     "+55" + customerDto.Phone,
		Cep:       customerDto.Cep,
		UpdatedAt: time.Now(),
	}

	err := service.Repository.Update(ctx, customer)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return exceptions.NewAppException(http.StatusNotFound, "Customer not found", nil)
		}
		return exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}
	return nil
}
