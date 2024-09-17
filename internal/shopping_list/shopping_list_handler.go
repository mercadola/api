package shoppinglist

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mercadola/api/pkg"
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
		r.Get("/{customerId}", h.FindByCustomerId)
	})
}

func (handler *ShoppingListHandler) FindByCustomerId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	customer_id := chi.URLParam(r, "customerId")
	resp, err := handler.Service.FindByCustomerId(r.Context(), customer_id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := pkg.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
