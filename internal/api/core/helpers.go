package core

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/mold/v4"
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func Validate[T any](r *http.Request, conform *mold.Transformer, validate *validator.Validate) (*T, Response) {
	req := new(T)
	err := json.UnmarshalRead(r.Body, req)
	if err != nil {
		return nil, Err(http.StatusBadRequest, err)
	}
	err = conform.Struct(r.Context(), req)
	if err != nil {
		return nil, Err(http.StatusBadRequest, err)
	}
	err = validate.StructCtx(r.Context(), req)
	if err != nil {
		if validationErrors, ok := errors.AsType[validator.ValidationErrors](err); ok {
			response := ValidationError{
				Message: validationErrors.Error(),
				Errors:  make(map[string]string, len(validationErrors)),
			}
			for _, validationError := range validationErrors {
				response.Errors[validationError.Field()] = validationError.Tag()
			}
			return nil, JSON(http.StatusBadRequest, response)
		}
		return nil, Err(http.StatusInternalServerError, fmt.Errorf("failed to match err to validations errors struct: %w", err))
	}
	return req, nil
}
