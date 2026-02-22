package api

import (
	"fmt"
	"net/http"

	"github.com/zagvozdeen/ola/internal/api/core"
	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Service) getServices(r *http.Request, user *models.User) core.Response {
	services, err := s.store.GetAllServices(r.Context())
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get services: %w", err))
	}
	return core.JSON(http.StatusOK, services)
}
