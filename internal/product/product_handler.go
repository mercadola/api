package product

import (
	"encoding/json"
	"net/http"

	"github.com/mercadola/api/pkg"
)

type ProductHandler struct {
	Service ProductService
}

func NewHandler(ps *ProductService) *ProductHandler {
	return &ProductHandler{
		Service: *ps,
	}
}

func (handler *ProductHandler) Find(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := FindQueryParams{}
	query.Ean = r.URL.Query().Get("ean")
	query.Ncm = r.URL.Query().Get("ncm")
	resp, err := handler.Service.Find(r.Context(), query)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := pkg.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
