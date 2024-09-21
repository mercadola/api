package shoppinglist

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/go-playground/validator/v10"
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
		r.Use(jwtauth.Authenticator)
		r.Post("/", h.Create)
		r.Get("/", h.FindByCustomer)
		r.Delete("/{shopping_list_id}", h.Delete)
		r.Patch("/{shopping_list_id}", h.UpdateName)
	})
}

func (h *ShoppingListHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, claims, _ := jwtauth.FromContext(r.Context())
	var shoppinglistDto ShoppingListCreateDto

	err := json.NewDecoder(r.Body).Decode(&shoppinglistDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := exceptions.NewAppException(http.StatusBadRequest, fmt.Sprintf("Error trying decode request => %s", err.Error()), nil)
		json.NewEncoder(w).Encode(error)
		return
	}
	if err := exceptions.ValidateException(validator.New(), shoppinglistDto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := exceptions.NewAppException(http.StatusBadRequest, err.Error(), nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	shoppingList, error := h.Service.Create(r.Context(), &shoppinglistDto, claims["sub"].(string))
	if error != nil {
		w.WriteHeader(error.StatusCode)
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shoppingList)
}

func (h *ShoppingListHandler) UpdateName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, claims, _ := jwtauth.FromContext(r.Context())
	shopping_list_id := chi.URLParam(r, "shopping_list_id")
	var shoppinglistDto ShoppingListUpdateDto

	err := json.NewDecoder(r.Body).Decode(&shoppinglistDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := exceptions.NewAppException(http.StatusBadRequest, fmt.Sprintf("Error trying decode request => %s", err.Error()), nil)
		json.NewEncoder(w).Encode(error)
		return
	}
	if err := exceptions.ValidateException(validator.New(), shoppinglistDto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := exceptions.NewAppException(http.StatusBadRequest, err.Error(), nil)
		json.NewEncoder(w).Encode(error)
		return
	}
	error := h.Service.UpdateName(r.Context(), shoppinglistDto.Name, claims["sub"].(string), shopping_list_id)
	if error != nil {
		w.WriteHeader(error.StatusCode)
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *ShoppingListHandler) FindByCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, claims, _ := jwtauth.FromContext(r.Context())
	resp, err := handler.Service.FindByCustomerId(r.Context(), claims["sub"].(string))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := exceptions.NewAppException(http.StatusNotFound, err.Error(), nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (handler *ShoppingListHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, claims, _ := jwtauth.FromContext(r.Context())
	shopping_list_id := chi.URLParam(r, "shopping_list_id")
	err := handler.Service.Delete(r.Context(), claims["sub"].(string), shopping_list_id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := exceptions.NewAppException(http.StatusNotFound, err.Error(), nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
