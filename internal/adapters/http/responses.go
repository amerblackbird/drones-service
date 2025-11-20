package http

import (
	"drones/internal/core/domain"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func ResponseWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ResponseWithErrorMessage(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func ResponseWithResouseNotFound(w http.ResponseWriter, resource string) {
	ResponseWithCustomError(w, http.StatusNotFound, *domain.NewDomainError(
		domain.ResourceNotFoundError,
		resource+" not found",
		nil,
	))
}

func ResponseWithCustomError(w http.ResponseWriter, status int, domainErr domain.DomainError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": domainErr.Message, "code": string(domainErr.Code)})
}

func ResponseWithValidationError(w http.ResponseWriter, status int, errors domain.ValidationErrors) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{"errors": errors})
}

func ResponseWithError(w http.ResponseWriter, err error) {
	if err == io.EOF {
		ResponseWithCustomError(w, http.StatusBadRequest, *domain.NewDomainError(
			domain.BodyIsRequiredError,
			"request body is required",
			nil,
		))
	} else if domainErr, ok := err.(*domain.DomainError); ok {
		ResponseWithCustomError(w, http.StatusBadRequest, *domainErr)
	} else if _, ok := err.(*validator.InvalidValidationError); ok {
		ResponseWithCustomError(w, http.StatusBadRequest, *domain.NewDomainError(
			domain.BodyIsRequiredError,
			"Invalid request body",
			nil,
		))

	} else {
		ResponseWithErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}
}
