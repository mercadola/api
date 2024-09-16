package customer

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mercadola/api/pkg"
)

type CustomerHandler struct {
	Service CustomerService
}

func NewHandler(cs *CustomerService) *CustomerHandler {
	return &CustomerHandler{
		Service: *cs,
	}
}

func (h *CustomerHandler) RegisterRoutes(r *chi.Mux) {
	r.Route("/customers", func(r chi.Router) {
		r.Get("/", h.Find)
	})
}

func (handler *CustomerHandler) Find(w http.ResponseWriter, r *http.Request) {
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
