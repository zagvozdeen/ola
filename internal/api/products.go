package api

import (
	"fmt"
	"net/http"

	"github.com/zagvozdeen/ola/internal/api/core"
	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Service) getProducts(r *http.Request, user *models.User) core.Response {
	products, err := s.store.GetAllProducts(r.Context())
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get products: %w", err))
	}
	return core.JSON(http.StatusOK, products)
}
