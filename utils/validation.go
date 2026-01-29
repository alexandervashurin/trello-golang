package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

// ValidationErrors представляет список ошибок валидации
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// HandleValidationError обрабатывает ошибку валидации и возвращает структурированный ответ
func HandleValidationError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var validationErrs validator.ValidationErrors
	if ok := err.(validator.ValidationErrors); ok != nil {
		validationErrs = ok
	} else {
		// Другая ошибка валидации
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	errors := make([]ValidationError, 0, len(validationErrs))

	for _, e := range validationErrs {
		field := e.Field()
		tag := e.Tag()
		value := ""

		if e.Param() != "" {
			value = e.Param()
		}

		errors = append(errors, ValidationError{
			Field: field,
			Tag:   tag,
			Value: value,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	response := ValidationErrors{Errors: errors}
	json.NewEncoder(w).Encode(response)
}

// RespondWithError отправляет ответ об ошибке
func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	type ErrorResponse struct {
		Error string `json:"error"`
	}

	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// RespondWithSuccess отправляет успешный ответ
func RespondWithSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
