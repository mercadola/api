package product

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ProductMock struct {
	mock.Mock
}

func (m *ProductMock) New(description, thumbnail, ncm, typePackaging, barcodeImage string, ean int64, width, height, length float64, quantityPackaging int) *Product {
	args := m.Called(description, thumbnail, ncm, typePackaging, barcodeImage, ean, width, height, length, quantityPackaging)
	return args.Get(0).(*Product)
}

type ProductRepositoryMock struct {
	mock.Mock
}

func (m *ProductRepositoryMock) Create(ctx context.Context, product *Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) Find(ctx context.Context, ean, ncm string) (*[]Product, error) {
	args := m.Called(ctx, ean, ncm)
	return args.Get(0).(*[]Product), args.Error(1)
}
