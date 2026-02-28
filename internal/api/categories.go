package api

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/zagvozdeen/ola/internal/api/core"
	"github.com/zagvozdeen/ola/internal/store/models"
)

var categorySlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:[-_][a-z0-9]+)*$`)

type upsertCategoryRequest struct {
	Slug string `json:"slug" mold:"trim,lcase" validate:"required,max=255"`
	Name string `json:"name" mold:"trim" validate:"required,max=255"`
}

func isValidCategorySlug(slug string) bool {
	return categorySlugPattern.MatchString(slug)
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

	slug := r.PathValue("slug")
	if !isValidCategorySlug(slug) {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid category slug"))
	}

	category, err := s.store.GetCategoryBySlug(r.Context(), slug)
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

	now := time.Now()
	category := &models.Category{
		Slug:      req.Slug,
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if !isValidCategorySlug(category.Slug) {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid category slug"))
	}

	err := s.store.CreateCategory(r.Context(), category)
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

	slug := r.PathValue("slug")
	if !isValidCategorySlug(slug) {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid category slug"))
	}

	req, res := core.Validate[upsertCategoryRequest](r, s.conform, s.validate)
	if res != nil {
		return res
	}
	if !isValidCategorySlug(req.Slug) {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid category slug"))
	}

	category, err := s.store.GetCategoryBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return core.Err(http.StatusNotFound, fmt.Errorf("category not found"))
		}
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get category: %w", err))
	}

	category.Slug = req.Slug
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

	slug := r.PathValue("slug")
	if !isValidCategorySlug(slug) {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid category slug"))
	}

	err := s.store.DeleteCategoryBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return core.Err(http.StatusNotFound, fmt.Errorf("category not found"))
		}
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to delete category: %w", err))
	}

	return core.JSON(http.StatusNoContent, nil)
}
