package shoppinglist

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
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

func (h *ShoppingListHandler) RegisterRoutes(r *chi.Mux, tokenAuth *jwtauth.JWTAuth) {
	r.Route("/shopping-list", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Post("/", h.Create)
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

func (h *ShoppingListHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var shoppinglistDto ShoppingListDto

	err := json.NewDecoder(r.Body).Decode(&shoppinglistDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := exceptions.NewAppException(http.StatusBadRequest, fmt.Sprintf("Error trying decode request => %s", err.Error()), nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	customer, err := h.Service.Create(r.Context(), &shoppinglistDto)
	if err != nil {
		if err, ok := err.(*exceptions.AppException); ok {
			w.WriteHeader(err.StatusCode)
			json.NewEncoder(w).Encode(err)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		error := exceptions.NewAppException(http.StatusInternalServerError, err.Error(), nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}
