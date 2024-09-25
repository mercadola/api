package test_mocks

import (
	"context"

	"github.com/mercadola/api/internal/customer"
	"github.com/mercadola/api/pkg/utils"
	"github.com/stretchr/testify/mock"
)

type CustomerMongoRepositoryMock struct {
	mock.Mock
	customer.CustomerRepositoryInterface
}

func NewCustomerMongoRepositoryMock() *CustomerMongoRepositoryMock {
	return &CustomerMongoRepositoryMock{}

}

func (m *CustomerMongoRepositoryMock) Create(ctx context.Context, customer *customer.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *CustomerMongoRepositoryMock) Delete(ctx context.Context, id string) (*utils.DeleteResult, error) {
	args := m.Called(ctx, id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*utils.DeleteResult), args.Error(1)
}

func (m *CustomerMongoRepositoryMock) Find(ctx context.Context, query *customer.FindQueryParams) (*[]customer.Customer, error) {
	args := m.Called(ctx, query)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]customer.Customer), args.Error(1)
}

func (m *CustomerMongoRepositoryMock) FindByEmail(ctx context.Context, findByEmailInput *customer.FindByEmailInput) (*customer.Customer, error) {
	args := m.Called(ctx, findByEmailInput)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*customer.Customer), args.Error(1)
}

func (m *CustomerMongoRepositoryMock) FindById(ctx context.Context, id string) (*customer.Customer, error) {
	args := m.Called(ctx, id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*customer.Customer), args.Error(1)
}

func (m *CustomerMongoRepositoryMock) InactiveCustomer(ctx context.Context, id string) (*utils.UpdateResult, error) {
	args := m.Called(ctx, id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*utils.UpdateResult), args.Error(1)
}

func (m *CustomerMongoRepositoryMock) PositivateCustomer(ctx context.Context, id string) (*utils.UpdateResult, error) {
	args := m.Called(ctx, id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*utils.UpdateResult), args.Error(1)
}

func (m *CustomerMongoRepositoryMock) Update(ctx context.Context, customer *customer.Customer) (*utils.UpdateResult, error) {
	args := m.Called(ctx, customer)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*utils.UpdateResult), args.Error(1)
}
