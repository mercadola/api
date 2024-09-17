package customer

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mercadola/api/internal/shared/utils/exceptions"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		r.Post("/", h.Create)
		r.Get("/", h.Find)
		r.Get("/{id}", h.FindById)
		r.Delete("/{id}", h.Delete)
	})
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var customerDto CustomerDto

	err := json.NewDecoder(r.Body).Decode(&customerDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := exceptions.NewAppException(http.StatusBadRequest, "Error trying decode request", err.Error(), nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	customer, err := h.Service.Create(r.Context(), &customerDto)
	if err != nil {
		if err, ok := err.(exceptions.AppException); ok {
			w.WriteHeader(err.StatusCode)
			json.NewEncoder(w).Encode(err)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		error := exceptions.NewAppException(http.StatusInternalServerError, "Internal Server Error", err.Error(), nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := exceptions.NewAppException(http.StatusBadRequest, "Error trying decode request", err.Error(), nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	err = h.Service.Delete(r.Context(), objectId)
	if err != nil {
		if err, ok := err.(exceptions.AppException); ok {
			w.WriteHeader(err.StatusCode)
			json.NewEncoder(w).Encode(err)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		error := exceptions.NewAppException(http.StatusInternalServerError, "Internal Server Error", err.Error(), nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(nil)
}

func (handler *CustomerHandler) Find(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := findQueryParams{}
	query.Name = r.URL.Query().Get("name")
	query.Email = r.URL.Query().Get("email")
	query.CPF = r.URL.Query().Get("cpf")
	resp, err := handler.Service.Find(r.Context(), query)
	if err != nil {
		if err, ok := err.(exceptions.AppException); ok {
			w.WriteHeader(err.StatusCode)
			json.NewEncoder(w).Encode(err)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		error := exceptions.NewAppException(http.StatusInternalServerError, "Internal Server Error", err.Error(), nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (handler *CustomerHandler) FindById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := exceptions.NewAppException(http.StatusBadRequest, "Error trying decode request", err.Error(), nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	customer, err := handler.Service.FindById(r.Context(), objectId)
	if err != nil {
		if err, ok := err.(exceptions.AppException); ok {
			w.WriteHeader(err.StatusCode)
			json.NewEncoder(w).Encode(err)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		error := exceptions.NewAppException(http.StatusInternalServerError, "Internal Server Error", err.Error(), nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}
