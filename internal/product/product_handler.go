package product

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mercadola/api/pkg/exceptions"
)

type ProductHandler struct {
	Service ProductService
}

func NewHandler(ps *ProductService) *ProductHandler {
	return &ProductHandler{
		Service: *ps,
	}
}

func (h *ProductHandler) RegisterRoutes(r *chi.Mux) {
	r.Route("/products", func(r chi.Router) {
		r.Get("/", h.Find)
	})
}

func (handler *ProductHandler) Find(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := FindProductQueryParams{}
	query.Ean = r.URL.Query().Get("ean")
	query.Ncm = r.URL.Query().Get("ncm")
	resp, err := handler.Service.Find(r.Context(), query)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := exceptions.NewAppException(http.StatusNotFound, err.Error(), nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
