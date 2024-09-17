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

func NewAppException(statusCode int, reason string, message string, validationDetails *[]ValidationException) *AppException {
	return newAppException(statusCode, reason, message, validationDetails)
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
