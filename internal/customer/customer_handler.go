package customer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/mercadola/api/pkg/exceptions"
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
	r.Post("/authenticate", h.Authenticate)
	r.Route("/customers", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Patch("/{id}/positivate", h.PositivateCustomer)
		r.Patch("/{id}/inactive", h.PositivateCustomer)
		r.Put("/{id}", h.Update)
		r.Get("/", h.Find)
		r.Get("/{id}", h.FindById)
		r.Delete("/{id}", h.Delete)
	})
}

func (h *CustomerHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jwtAuth := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("jwtExpiresIn").(int)

	var authenticateInput AuthenticateInput

	err := json.NewDecoder(r.Body).Decode(&authenticateInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := exceptions.NewAppException(http.StatusBadRequest, fmt.Sprintf("Error trying decode request => %s", err.Error()), nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	if err = authenticateInput.Validate(); err != nil {
		handleCustomError(w, err)
		return
	}

	customer, err := h.Service.Authenticate(r.Context(), authenticateInput)
	if err != nil {
		handleCustomError(w, err)
		return
	}

	_, token, _ := jwtAuth.Encode(map[string]interface{}{
		"sub": customer.ID,
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})

	accessToken := AutenticateOutput{AccessToken: token}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var customerDto CustomerDto

	err := json.NewDecoder(r.Body).Decode(&customerDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := exceptions.NewAppException(http.StatusBadRequest, fmt.Sprintf("Error trying decode request => %s", err.Error()), nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	if err = customerDto.Validate(); err != nil {
		handleCustomError(w, err)
		return
	}

	customer, err := h.Service.Create(r.Context(), &customerDto)
	if err != nil {
		handleCustomError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := exceptions.NewAppException(http.StatusBadRequest, "Invalid id", nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	err := h.Service.Delete(r.Context(), id)
	if err != nil {
		handleCustomError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(nil)
}

func (handler *CustomerHandler) Find(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := FindQueryParams{}
	query.Email = r.URL.Query().Get("email")
	query.CPF = r.URL.Query().Get("cpf")

	if err := query.Validate(); err != nil {
		handleCustomError(w, err)
		return
	}

	customers, err := handler.Service.Find(r.Context(), query)
	if err != nil {
		handleCustomError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customers)
}

func (handler *CustomerHandler) FindById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := exceptions.NewAppException(http.StatusBadRequest, "Invalid id", nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	customer, err := handler.Service.FindById(r.Context(), id)
	if err != nil {
		handleCustomError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) PositivateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := exceptions.NewAppException(http.StatusBadRequest, "Invalid id", nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	err := h.Service.PositivateCustomer(r.Context(), id)
	if err != nil {
		handleCustomError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CustomerHandler) InactiveCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := exceptions.NewAppException(http.StatusBadRequest, "Invalid id", nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	err := h.Service.InactiveCustomer(r.Context(), id)
	if err != nil {
		handleCustomError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := exceptions.NewAppException(http.StatusBadRequest, "Invalid id", nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	var customerDto CustomerDto

	err := json.NewDecoder(r.Body).Decode(&customerDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := exceptions.NewAppException(http.StatusBadRequest, fmt.Sprintf("Error trying decode request => %s", err.Error()), nil)
		json.NewEncoder(w).Encode(error)
		return
	}

	if err = customerDto.Validate(); err != nil {
		handleCustomError(w, err)
		return
	}

	err = h.Service.Update(r.Context(), id, &customerDto)
	if err != nil {
		handleCustomError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func buildObjectId(id string) (primitive.ObjectID, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, exceptions.NewAppException(http.StatusBadRequest, fmt.Sprintf("Error trying decode request => %s", err.Error()), nil)
	}
	return objectId, nil
}

func handleCustomError(w http.ResponseWriter, err error) {
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
