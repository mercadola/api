package test_mocks

import (
	"github.com/mercadola/api/internal/customer"
	"github.com/stretchr/testify/mock"
)

type CustomerMock struct {
	mock.Mock
	customer.CustomerInterface
}

func NewCustomerMock() *CustomerMock {
	return &CustomerMock{}
}

func (m *CustomerMock) New(customerDto *customer.CustomerDto) (*customer.Customer, error) {
	args := m.Called(customerDto)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*customer.Customer), args.Error(1)
}
