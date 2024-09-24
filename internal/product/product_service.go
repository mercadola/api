package product

import (
	"context"
	"net/http"

	"github.com/mercadola/api/pkg/exceptions"
)

type ProductService struct {
	Repository ProductRepositoryInterface
}

func NewService(pr ProductRepositoryInterface) *ProductService {
	return &ProductService{
		Repository: pr,
	}
}
func (service *ProductService) Find(ctx context.Context, ean, ncm string) (*[]Product, *exceptions.AppException) {
	products, err := service.Repository.Find(ctx, ean, ncm)
	if err != nil {
		return nil, exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}

	return products, nil
}

func (service *ProductService) FindById(ctx context.Context, id string) (*Product, *exceptions.AppException) {
	product, err := service.Repository.FindById(ctx, id)
	if err != nil {
		return nil, exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}
	if product == nil {
		return nil, exceptions.NewAppException(http.StatusNotFound, "Nenhum documento encontrado", nil)
	}

	return product, nil
}

func (service *ProductService) CreateByEan(ctx context.Context, ean string) (*Product, *exceptions.AppException) {
	products, err := service.Repository.Find(ctx, ean, "")
	if err != nil {
		return nil, exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}
	if len(*products) > 0 {
		return nil, exceptions.NewAppException(http.StatusConflict, "Produto já cadastrado", nil)
	}
	product, err := GetCosmosProductByEan(ean)
	if err != nil {
		return nil, exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
	}
	service.Repository.Create(ctx, product)

	return nil, nil
}
