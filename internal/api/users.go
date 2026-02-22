package api

import (
	"fmt"
	"net/http"

	"github.com/zagvozdeen/ola/internal/api/core"
	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Service) getMe(r *http.Request, user *models.User) core.Response {
	return core.JSON(http.StatusOK, user)
}

func (s *Service) getUsers(r *http.Request, user *models.User) core.Response {
	res := allowForAdmin(user)
	if res != nil {
		return res
	}

	users, err := s.store.GetAllUsers(r.Context())
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get users: %w", err))
	}
	return core.JSON(http.StatusOK, users)
}
