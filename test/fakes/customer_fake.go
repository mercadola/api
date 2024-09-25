package test_fakes

import (
	"time"

	"github.com/google/uuid"
	"github.com/mercadola/api/internal/customer"
)

var firstObjectId = uuid.New().String()
var lastObjectId = uuid.New().String()
var FakeCustomerObjectId = "86e7d3f9-e00a-48c5-929d-12186818f486"

var FakeCustomerMap = map[string]*customer.Customer{
	firstObjectId: {
		ID:    firstObjectId,
		Name:  "any_name",
		Email: "any_email2@test.com",
		// Password:  "$2a$10$xD0Tn7XMGVzJN0tXltprf.JJ.eJFHi7U3KTAJtq51ud8bhJ4cQaz6",
		Password:  "$2a$10$vhoTn67.ELugKBsO2TCa8e8.7rqxvqPfeqdifd32g42Hos2hi4cOq",
		CPF:       "1234567890",
		Phone:     "+5521999999999",
		Cep:       "12345123",
		Gender:    customer.Male,
		Birthday:  time.Date(1990, time.October, 9, 18, 32, 0, 0, time.UTC),
		Active:    true,
		CreatedAt: time.Date(2024, time.September, 23, 18, 32, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, time.September, 23, 18, 32, 0, 0, time.UTC),
	},
	FakeCustomerObjectId: {
		ID:    FakeCustomerObjectId,
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
	},
	lastObjectId: {
		ID:    lastObjectId,
		Name:  "other_name",
		Email: "other_email@test.com",
		// Password:  "$2a$10$xD0Tn7XMGVzJN0tXltprf.JJ.eJFHi7U3KTAJtq51ud8bhJ4cQaz6",
		Password:  "$2a$10$vhoTn67.ELugKBsO2TCa8e8.7rqxvqPfeqdifd32g42Hos2hi4cOq",
		CPF:       "1234567892",
		Phone:     "+5521999999999",
		Cep:       "12345123",
		Gender:    customer.Female,
		Birthday:  time.Date(1990, time.October, 9, 18, 32, 0, 0, time.UTC),
		Active:    true,
		CreatedAt: time.Date(2024, time.September, 23, 18, 32, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, time.September, 23, 18, 32, 0, 0, time.UTC),
	},
}

func GetFakeCustomerSuccess() *customer.Customer {
	return FakeCustomerMap[FakeCustomerObjectId]
}
