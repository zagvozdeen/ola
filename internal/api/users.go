package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/api/core"
	"github.com/zagvozdeen/ola/internal/store/enums"
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

func (s *Service) getUser(r *http.Request, user *models.User) core.Response {
	res := allowForAdmin(user)
	if res != nil {
		return res
	}

	uid, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid user uuid: %w", err))
	}

	target, err := s.store.GetUserByUUID(r.Context(), uid)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return core.Err(http.StatusNotFound, fmt.Errorf("user not found"))
		}
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get user: %w", err))
	}

	return core.JSON(http.StatusOK, target)
}

type updateUserRoleRequest struct {
	Role string `json:"role" mold:"trim,lcase" validate:"required,oneof=user moderator admin"`
}

func (s *Service) updateUserRole(r *http.Request, user *models.User) core.Response {
	res := allowForAdmin(user)
	if res != nil {
		return res
	}

	uid, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid user uuid: %w", err))
	}

	req, res := core.Validate[updateUserRoleRequest](r, s.conform, s.validate)
	if res != nil {
		return res
	}

	role, err := enums.NewUserRole(req.Role)
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid user role: %w", err))
	}

	target, err := s.store.GetUserByUUID(r.Context(), uid)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return core.Err(http.StatusNotFound, fmt.Errorf("user not found"))
		}
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get user: %w", err))
	}

	if target.ID == user.ID {
		return core.Err(http.StatusBadRequest, fmt.Errorf("you can not change your own role"))
	}

	if target.Role == enums.UserRoleAdmin && role != enums.UserRoleAdmin {
		adminCount, countErr := s.store.CountUsersByRole(r.Context(), enums.UserRoleAdmin)
		if countErr != nil {
			return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to count admins: %w", countErr))
		}
		if adminCount <= 1 {
			return core.Err(http.StatusBadRequest, fmt.Errorf("you can not demote the last admin"))
		}
	}

	err = s.store.UpdateUserRole(r.Context(), target.ID, role)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to update user role: %w", err))
	}

	updated, err := s.store.GetUserByUUID(r.Context(), uid)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to load updated user: %w", err))
	}

	return core.JSON(http.StatusOK, updated)
}
