package errorhandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"zentinel/internal/domain"
	"zentinel/internal/infrastructure/driving/http/dto/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func HandleRequestError(c *gin.Context, err error) {
	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		msg := fmt.Sprintf("Field '%s' must be of type %s", unmarshalTypeError.Field, unmarshalTypeError.Type.String())
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_TYPE", "Invalid data type", msg))
		return
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var errorMessages []string
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' is required", e.Field()))
			case "email":
				errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' must be a valid email", e.Field()))
			case "uuid":
				errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' must be a valid UUID", e.Field()))
			default:
				errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' failed validation: %s", e.Field(), e.Tag()))
			}
		}
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "Validation failed", strings.Join(errorMessages, "; ")))
		return
	}

	c.JSON(http.StatusBadRequest, response.NewErrorResponse("BAD_REQUEST", "Malformed JSON or invalid request", err.Error()))
}

func HandleDomainError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, response.NewErrorResponse("NOT_FOUND", "Resource not found", err.Error()))
	case errors.Is(err, domain.ErrConflict):
		c.JSON(http.StatusConflict, response.NewErrorResponse("CONFLICT", "Resource conflict", err.Error()))
	case errors.Is(err, domain.ErrDocumentAlreadyExists):
		c.JSON(http.StatusConflict, response.NewErrorResponse("DOCUMENT_EXISTS", "Document already exists", err.Error()))
	case errors.Is(err, domain.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_INPUT", "Invalid input", err.Error()))
	case errors.Is(err, domain.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("UNAUTHORIZED", "Unauthorized", err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "An unexpected error occurred", err.Error()))
	}
}
