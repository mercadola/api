package product

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/go-playground/validator/v10"
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

func (h *ProductHandler) RegisterRoutes(r *chi.Mux, tokenAuth *jwtauth.JWTAuth) {
	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/", h.Find)
		r.Get("/{product_id}", h.FindById)
		r.Post("/ean", h.CreateByEan)
	})
}
func (h *ProductHandler) CreateByEan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var createByEanDto CreateByEanDto

	err := json.NewDecoder(r.Body).Decode(&createByEanDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := exceptions.NewAppException(http.StatusBadRequest, fmt.Sprintf("Error trying decode request => %s", err.Error()), nil)
		json.NewEncoder(w).Encode(error)
		return
	}
	if err := exceptions.ValidateException(validator.New(), createByEanDto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := exceptions.NewAppException(http.StatusBadRequest, err.Error(), nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	shoppingList, error := h.Service.CreateByEan(r.Context(), createByEanDto.Ean)
	if error != nil {
		w.WriteHeader(error.StatusCode)
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shoppingList)
}

func (handler *ProductHandler) Find(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := FindProductQueryParams{
		Ean: r.URL.Query().Get("ean"),
		Ncm: r.URL.Query().Get("ncm"),
	}
	resp, err := handler.Service.Find(r.Context(), query.Ean, query.Ncm)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (handler *ProductHandler) FindById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	product_id := chi.URLParam(r, "product_id")
	resp, err := handler.Service.FindById(r.Context(), product_id)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
