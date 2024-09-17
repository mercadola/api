package exceptions

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AppException struct {
	StatusCode        int                    `json:"statusCode"`
	Reason            string                 `json:"reason"`
	Message           string                 `json:"message"`
	ValidationDetails *[]ValidationException `json:"validationDetails"`
}

func (e AppException) Error() string {
	return e.Message
}

type ValidationException struct {
	Field    string `json:"field"`
	Tag      string `json:"tag"`
	Received any    `json:"received"`
}

func newAppException(statusCode int, reason string, message string, validationDetails *[]ValidationException) *AppException {
	return &AppException{
		StatusCode:        statusCode,
		Reason:            reason,
		Message:           message,
		ValidationDetails: validationDetails,
	}
}

func NewAppException(statusCode int, message string, validationDetails *[]ValidationException) *AppException {
	switch statusCode {
	case http.StatusBadRequest:
		return newAppException(statusCode, "Bad Request", message, validationDetails)
	case http.StatusUnauthorized:
		return newAppException(statusCode, "Unauthorized", message, validationDetails)
	case http.StatusForbidden:
		return newAppException(statusCode, "Forbidden", message, validationDetails)
	case http.StatusNotFound:
		return newAppException(statusCode, "Not Found", message, validationDetails)
	case http.StatusConflict:
		return newAppException(statusCode, "Conflict", message, validationDetails)
	default:
		return newAppException(http.StatusInternalServerError, "Internal Server Error", message, validationDetails)
	}

}

func ValidateException(validate *validator.Validate, dto any) error {
	var errors []ValidationException
	log := slog.Default()

	err := validate.Struct(dto)

	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.Error("Error trying to validate struct", err)
		}

		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationException{
				Field:    err.Field(),
				Tag:      err.Tag(),
				Received: err.Value(),
			})
		}
		return newAppException(http.StatusBadRequest, "Bad Request", "Invalid request", &errors)
	}
	return nil
}
