package errors

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func RespondError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

func RespondValidationError(w http.ResponseWriter, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		fields := make(map[string][]string)

		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			if fields[field] == nil {
				fields[field] = []string{}
			}
			fields[field] = append(fields[field], getValidationMessage(fieldError))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)

		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error: ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: "Validation failed",
				Details: map[string]interface{}{
					"fields": fields,
				},
			},
		})
		return
	}

	RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
}

func getValidationMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return fieldError.Field() + " is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fieldError.Field() + " must be at least " + fieldError.Param() + " characters"
	case "max":
		return fieldError.Field() + " must be at most " + fieldError.Param() + " characters"
	default:
		return fieldError.Field() + " is invalid"
	}
}
