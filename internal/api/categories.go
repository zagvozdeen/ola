package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/api/core"
	"github.com/zagvozdeen/ola/internal/store/models"
)

type upsertCategoryRequest struct {
	Name string `json:"name" mold:"trim" validate:"required,max=255"`
}

func (s *Service) getCategories(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	categories, err := s.store.GetAllCategories(r.Context())
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get categories: %w", err))
	}
	return core.JSON(http.StatusOK, categories)
}

func (s *Service) getCategory(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	uid, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid category uuid: %w", err))
	}

	category, err := s.store.GetCategoryByUUID(r.Context(), uid)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return core.Err(http.StatusNotFound, fmt.Errorf("category not found"))
		}
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get category: %w", err))
	}

	return core.JSON(http.StatusOK, category)
}

func (s *Service) createCategory(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	req, res := core.Validate[upsertCategoryRequest](r, s.conform, s.validate)
	if res != nil {
		return res
	}

	uid, err := uuid.NewV7()
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to generate uuid v7: %w", err))
	}

	now := time.Now()
	category := &models.Category{
		UUID:      uid,
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.store.CreateCategory(r.Context(), category)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to create category: %w", err))
	}

	return core.JSON(http.StatusCreated, category)
}

func (s *Service) updateCategory(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	uid, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid category uuid: %w", err))
	}

	req, res := core.Validate[upsertCategoryRequest](r, s.conform, s.validate)
	if res != nil {
		return res
	}

	category, err := s.store.GetCategoryByUUID(r.Context(), uid)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return core.Err(http.StatusNotFound, fmt.Errorf("category not found"))
		}
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get category: %w", err))
	}

	category.Name = req.Name
	category.UpdatedAt = time.Now()

	err = s.store.UpdateCategory(r.Context(), category)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to update category: %w", err))
	}

	return core.JSON(http.StatusOK, category)
}

func (s *Service) deleteCategory(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	uid, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid category uuid: %w", err))
	}

	err = s.store.DeleteCategoryByUUID(r.Context(), uid)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return core.Err(http.StatusNotFound, fmt.Errorf("category not found"))
		}
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to delete category: %w", err))
	}

	return core.JSON(http.StatusNoContent, nil)
}
