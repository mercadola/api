package shoppinglist

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mercadola/api/pkg/exceptions"
)

type ShoppingListHandler struct {
	Service ShoppingListService
}

func NewHandler(ps *ShoppingListService) *ShoppingListHandler {
	return &ShoppingListHandler{
		Service: *ps,
	}
}

func (h *ShoppingListHandler) RegisterRoutes(r *chi.Mux) {
	r.Route("/shopping-list", func(r chi.Router) {
		r.Get("/{customer_id}", h.FindByCustomerId)
	})
}

func (handler *ShoppingListHandler) FindByCustomerId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	customer_id := chi.URLParam(r, "customer_id")
	resp, err := handler.Service.FindByCustomerId(r.Context(), customer_id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := exceptions.NewAppException(http.StatusNotFound, err.Error(), nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
