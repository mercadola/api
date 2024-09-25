package customer

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mercadola/api/pkg/exceptions"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerServiceInterface interface {
	Authenticate(ctx context.Context, authenticateInput AuthenticateInput) (*Customer, *exceptions.AppException)
	Create(ctx context.Context, customerDto *CustomerDto) (*Customer, *exceptions.AppException)
	Find(ctx context.Context, query FindQueryParams) (*[]Customer, *exceptions.AppException)
	FindById(ctx context.Context, id primitive.ObjectID) (*Customer, *exceptions.AppException)
	FindByEmail(ctx context.Context, findByEmailInput *FindByEmailInput) (*Customer, *exceptions.AppException)
	InactiveCustomer(ctx context.Context, id primitive.ObjectID) *exceptions.AppException
	PositivateCustomer(ctx context.Context, id primitive.ObjectID) *exceptions.AppException
	Update(ctx context.Context, id primitive.ObjectID, customerDto *CustomerDto) *exceptions.AppException
}

type CustomerService struct {
	Customer   CustomerInterface
	Logger     *slog.Logger
	Repository CustomerRepositoryInterface
}

func NewService(cr CustomerRepositoryInterface, logger *slog.Logger, customer CustomerInterface) *CustomerService {
	return &CustomerService{
		Repository: cr,
		Logger:     logger,
		Customer:   customer,
	}
}

func (service *CustomerService) Authenticate(ctx context.Context, authenticateInput AuthenticateInput) (*Customer, *exceptions.AppException) {
	finddedCustomer, err := service.FindByEmail(ctx, &FindByEmailInput{Email: authenticateInput.Email})

	if err != nil {
		return nil, exceptions.NewAppException(http.StatusUnauthorized, "Invalid credentials", nil)
	}

	if !finddedCustomer.validatePassword(authenticateInput.Password) {
		return nil, exceptions.NewAppException(http.StatusUnauthorized, "Invalid credentials", nil)
	}

	return finddedCustomer, nil
}

func (service *CustomerService) Create(ctx context.Context, customerDto *CustomerDto) (*Customer, *exceptions.AppException) {
	finddedCustomers, _ := service.Find(ctx, FindQueryParams{Email: customerDto.Email, CPF: customerDto.CPF})

	if finddedCustomers != nil && len(*finddedCustomers) > 0 {
		return nil, exceptions.NewAppException(http.StatusConflict, "Customer already exists", nil)
	}

	customer, err := service.Customer.New(customerDto)
	if err != nil {
		service.Logger.Error(fmt.Sprintf("Error trying create instance customer => %s", err.Error()))
		return nil, exceptions.NewAppException(http.StatusInternalServerError, "Error trying to create customer", nil)
	}

	err = service.Repository.Create(ctx, customer)
	if err != nil {
		service.Logger.Error(fmt.Sprintf("Error trying create customer in database => %s", err.Error()))
		return nil, exceptions.NewAppException(http.StatusInternalServerError, "Error trying to create customer", nil)
	}

	// TODO DISPARO DE E-MAIL DE BOAS VINDAS

	return customer, nil
}

func (service *CustomerService) Delete(ctx context.Context, id string) *exceptions.AppException {
	result, err := service.Repository.Delete(ctx, id)
	if err != nil {
		return exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}
	if result.DeletedCount == 0 {
		return exceptions.NewAppException(http.StatusNotFound, "no documents in result", nil)
	}
	return nil
}

func (service *CustomerService) Find(ctx context.Context, query FindQueryParams) (*[]Customer, *exceptions.AppException) {
	customers, err := service.Repository.Find(ctx, &query)
	if err != nil {
		return nil, exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}
	return customers, nil
}

func (service *CustomerService) FindByEmail(ctx context.Context, findByEmail *FindByEmailInput) (*Customer, *exceptions.AppException) {
	customer, err := service.Repository.FindByEmail(ctx, findByEmail)
	if err != nil {
		return nil, exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}

	if customer == nil {
		return nil, exceptions.NewAppException(http.StatusNotFound, "Customer not found", nil)
	}

	return customer, nil
}

func (service *CustomerService) FindById(ctx context.Context, id string) (*Customer, *exceptions.AppException) {
	customer, err := service.Repository.FindById(ctx, id)
	if err != nil {
		return nil, exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}

	if customer == nil {
		return nil, exceptions.NewAppException(http.StatusNotFound, "Customer not found", nil)
	}

	return customer, nil
}

func (service *CustomerService) InactiveCustomer(ctx context.Context, id string) *exceptions.AppException {
	result, err := service.Repository.InactiveCustomer(ctx, id)
	if err != nil {
		return exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}

	if result.ModifiedCount == 0 {
		return exceptions.NewAppException(http.StatusNotFound, "no documents in result", nil)
	}

	return nil
}

func (service *CustomerService) PositivateCustomer(ctx context.Context, id string) *exceptions.AppException {
	result, err := service.Repository.PositivateCustomer(ctx, id)
	if err != nil {
		return exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}

	if result.ModifiedCount == 0 {
		return exceptions.NewAppException(http.StatusNotFound, "no documents in result", nil)
	}

	return nil
}

func (service *CustomerService) Update(ctx context.Context, id string, customerDto *CustomerDto) *exceptions.AppException {
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

	result, err := service.Repository.Update(ctx, &customer)

	if err != nil {
		return exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}

	if result.ModifiedCount == 0 {
		return exceptions.NewAppException(http.StatusNotFound, "no documents in result", nil)
	}

	return nil
}
