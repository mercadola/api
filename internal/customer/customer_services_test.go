package customer_test

import (
	"context"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/mercadola/api/internal/customer"
	"github.com/mercadola/api/pkg/exceptions"
	"github.com/mercadola/api/test"
	test_fakes "github.com/mercadola/api/test/fakes"
	test_mocks "github.com/mercadola/api/test/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_AuthenticateCustomer(t *testing.T) {
	expectedErr := exceptions.NewAppException(http.StatusUnauthorized, "Invalid credentials", nil)

	testCases := []test.TestCase{
		{
			Name: "should authenticate a customer",
			InputParams: customer.AuthenticateInput{
				Email:    "any_email@test.com",
				Password: "123456789",
			},
			Case: test.Success,
			FunctionsCalledMock: []test.CaseFunction{
				{
					Name:                "FindByEmail",
					InputParams:         &customer.FindByEmailInput{Email: "any_email@test.com"},
					Expected:            test_fakes.GetFakeCustomerSuccess(),
					ExpectedErr:         nil,
					AssertNumberOfCalls: 1,
				},
			},
			Expected:    test_fakes.GetFakeCustomerSuccess(),
			ExpectedErr: nil,
		},
		{
			Name: "should return error when password is wrong",
			InputParams: customer.AuthenticateInput{
				Email:    "any_email@test.com",
				Password: "wrong_password",
			},
			Case: test.Failure,
			FunctionsCalledMock: []test.CaseFunction{
				{
					Name:                "FindByEmail",
					InputParams:         &customer.FindByEmailInput{Email: "any_email@test.com"},
					Expected:            test_fakes.GetFakeCustomerSuccess(),
					ExpectedErr:         nil,
					AssertNumberOfCalls: 1,
				},
			},
			Expected:    nil,
			ExpectedErr: expectedErr,
		},
		{
			Name: "should return error when email is wrong",
			InputParams: customer.AuthenticateInput{
				Email:    "wrong_email@test.com",
				Password: "any_password",
			},
			Case: test.Failure,
			FunctionsCalledMock: []test.CaseFunction{
				{
					Name:                "FindByEmail",
					InputParams:         &customer.FindByEmailInput{Email: "wrong_email@test.com"},
					Expected:            nil,
					ExpectedErr:         expectedErr,
					AssertNumberOfCalls: 1,
				},
			},
			Expected:    nil,
			ExpectedErr: expectedErr,
		},
	}

	logger := slog.Default()

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			mockRepository := test_mocks.NewCustomerMongoRepositoryMock()
			mockCustomer := test_mocks.NewCustomerMock()

			customerService := customer.NewService(mockRepository, logger, mockCustomer)

			for _, function := range tc.FunctionsCalledMock {
				mockRepository.On(function.Name, ctx, function.InputParams).Return(function.Expected, function.ExpectedErr)
			}

			result, err := customerService.Authenticate(ctx, tc.InputParams.(customer.AuthenticateInput))

			if tc.Case == test.Success {
				assert.Nil(t, err)
				assert.Equal(t, tc.Expected.(*customer.Customer), result)
			}

			if tc.Case == test.Failure {
				assert.Nil(t, result)
				assert.Equal(t, tc.Expected, nil)
				assert.NotNil(t, err)
				assert.Equal(t, tc.ExpectedErr, err)
			}

			for _, function := range tc.FunctionsCalledMock {
				if function.AssertNumberOfCalls > 0 {
					mockRepository.AssertNumberOfCalls(t, function.Name, function.AssertNumberOfCalls)
				} else {
					mockRepository.AssertNotCalled(t, function.Name)
				}
			}
		})
	}
}

func Test_CreateCustomer(t *testing.T) {
	// conflictErr := exceptions.NewAppException(http.StatusConflict, "Customer already exists", nil)
	// genericErr := exceptions.NewAppException(http.StatusInternalServerError, "Error trying to create customer", nil)
	customerDtoCaseSuccess := customer.CustomerDto{
		Name:     "any_name",
		Email:    "any_email@test.com",
		Password: "123456789",
		CPF:      "1234567891",
		Phone:    "21999999999",
		Cep:      "12345123",
		Gender:   customer.Male,
		Birthday: time.Date(1990, time.October, 9, 18, 32, 0, 0, time.UTC),
	}

	testCases := []test.TestCase{
		{
			Name:        "should create a customer",
			InputParams: &customerDtoCaseSuccess,
			Case:        test.Success,
			FunctionsCalledMock: []test.CaseFunction{
				{
					Name:                "Find",
					MockName:            "repository",
					InputParams:         &customer.FindQueryParams{Email: customerDtoCaseSuccess.Email, CPF: customerDtoCaseSuccess.CPF},
					Expected:            &[]customer.Customer{},
					ExpectedErr:         nil,
					AssertNumberOfCalls: 1,
				},
				{
					Name:                "New",
					MockName:            "customer",
					InputParams:         &customerDtoCaseSuccess,
					Expected:            test_fakes.GetFakeCustomerSuccess(),
					ExpectedErr:         nil,
					AssertNumberOfCalls: 1,
				},
				{
					Name:                "Create",
					MockName:            "repository",
					InputParams:         test_fakes.GetFakeCustomerSuccess(),
					Expected:            nil,
					ExpectedErr:         nil,
					AssertNumberOfCalls: 1,
				},
			},
			Expected:    test_fakes.GetFakeCustomerSuccess(),
			ExpectedErr: nil,
		},
		// {
		// 	Name: "should return conflict error when email is already in use",
		// 	InputParams: customer.CustomerDto{
		// 		Name:     "any_name",
		// 		Email:    "any_email@test.com",
		// 		Password: "123456789",
		// 		CPF:      "1234567891",
		// 		Phone:    "+5521999999999",
		// 		Cep:      "12345123",
		// 		Gender:   customer.Male,
		// 		Birthday: time.Date(1990, time.October, 9, 18, 32, 0, 0, time.UTC),
		// 	},
		// 	Case: test.Failure,
		// 	FunctionsCalledMock: []test.CaseFunction{
		// 		{
		// 			Name:                "Find",
		// 			MockName:            "repository",
		// 			InputParams:         &customer.FindQueryParams{Email: "any_email@test.com", CPF: "1234567891"},
		// 			Expected:            test_fakes.GetFakeCustomerByIdSuccess(),
		// 			ExpectedErr:         nil,
		// 			AssertNumberOfCalls: 1,
		// 		},
		// 		{
		// 			Name:                "New",
		// 			MockName:            "customer",
		// 			AssertNumberOfCalls: 0,
		// 		},
		// 		{
		// 			Name:                "Create",
		// 			MockName:            "repository",
		// 			AssertNumberOfCalls: 0,
		// 		},
		// 	},
		// 	Expected:    test_fakes.GetFakeCustomerByIdSuccess(),
		// 	ExpectedErr: nil,
		// },
	}

	logger := slog.Default()

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()

			mockRepository := test_mocks.NewCustomerMongoRepositoryMock()
			mockCustomer := test_mocks.NewCustomerMock()

			expect := &customer.Customer{
				ID:    test_fakes.FakeCustomerObjectId,
				Name:  "any_name",
				Email: "any_email@test.com",
				// Password:  "$2a$10$xD0Tn7XMGVzJN0tXltprf.JJ.eJFHi7U3KTAJtq51ud8bhJ4cQaz6",
				Password:  "$2a$10$vhoTn67.ELugKBsO2TCa8e8.7rqxvqPfeqdifd32g42Hos2hi4cOq",
				CPF:       "1234567891",
				Phone:     "+5521999999999",
				Cep:       "12345123",
				Gender:    customer.Male,
				Birthday:  time.Date(1990, time.October, 9, 18, 32, 0, 0, time.UTC),
				Active:    true,
				CreatedAt: time.Date(2024, time.September, 23, 18, 32, 0, 0, time.UTC),
				UpdatedAt: time.Date(2024, time.September, 23, 18, 32, 0, 0, time.UTC),
			}

			mockCustomer.On("New", &customerDtoCaseSuccess).Return(expect, nil).Once()

			for _, function := range tc.FunctionsCalledMock {
				if function.AssertNumberOfCalls > 0 {
					// if function.MockName == "customer" {
					// 	mockCustomer.On(function.Name, function.InputParams.(*customer.CustomerDto)).Return(function.Expected.(*customer.Customer), function.ExpectedErr)
					// }
					if function.MockName == "repository" {
						mockRepository.On(function.Name, ctx, function.InputParams).Return(function.Expected, function.ExpectedErr)
					}
				}
			}

			customerService := customer.NewService(mockRepository, logger, mockCustomer)

			result, err := customerService.Create(ctx, tc.InputParams.(*customer.CustomerDto))

			if tc.Case == test.Success {
				assert.Nil(t, err)
				assert.Equal(t, tc.Expected.(*customer.Customer), result)
			}

			if tc.Case == test.Failure {
				assert.Nil(t, result)
				assert.Equal(t, tc.Expected, nil)
				assert.NotNil(t, err)
				assert.Equal(t, tc.ExpectedErr, err)
			}

			for _, function := range tc.FunctionsCalledMock {
				if function.AssertNumberOfCalls > 0 {
					if function.MockName == "mockCustomer" {
						mockCustomer.AssertNumberOfCalls(t, function.Name, function.AssertNumberOfCalls)
					}
					if function.MockName == "mockRepository" {
						mockRepository.AssertNumberOfCalls(t, function.Name, function.AssertNumberOfCalls)
					}
				} else {
					if function.MockName == "mockCustomer" {
						mockCustomer.AssertNotCalled(t, function.Name)
					}
					if function.MockName == "mockRepository" {
						mockRepository.AssertNotCalled(t, function.Name)
					}
				}
			}
		})
	}
}
