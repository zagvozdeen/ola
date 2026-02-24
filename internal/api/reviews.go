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

type upsertReviewRequest struct {
	Name        string `json:"name" mold:"trim" validate:"required,max=255"`
	Content     string `json:"content" mold:"trim" validate:"required,max=3000"`
	FileID      int    `json:"file_id" validate:"required,gt=0"`
	PublishedAt string `json:"published_at" mold:"trim" validate:"required"`
}

func (s *Service) getReviews(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	reviews, err := s.store.GetAllReviews(r.Context())
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get reviews: %w", err))
	}
	return core.JSON(http.StatusOK, reviews)
}

func (s *Service) getReview(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	uid, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid review uuid: %w", err))
	}

	review, err := s.store.GetReviewByUUID(r.Context(), uid)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return core.Err(http.StatusNotFound, fmt.Errorf("review not found"))
		}
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get review: %w", err))
	}

	return core.JSON(http.StatusOK, review)
}

func (s *Service) createReview(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	req, res := core.Validate[upsertReviewRequest](r, s.conform, s.validate)
	if res != nil {
		return res
	}

	publishedAt, err := time.Parse(time.RFC3339, req.PublishedAt)
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid published_at: %w", err))
	}

	uid, err := uuid.NewV7()
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to generate uuid v7: %w", err))
	}

	now := time.Now()
	review := &models.Review{
		UUID:        uid,
		Name:        req.Name,
		Content:     req.Content,
		FileID:      req.FileID,
		UserID:      user.ID,
		PublishedAt: publishedAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = s.store.CreateReview(r.Context(), review)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to create review: %w", err))
	}

	created, err := s.store.GetReviewByUUID(r.Context(), uid)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to load created review: %w", err))
	}

	return core.JSON(http.StatusCreated, created)
}

func (s *Service) updateReview(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	uid, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid review uuid: %w", err))
	}

	req, res := core.Validate[upsertReviewRequest](r, s.conform, s.validate)
	if res != nil {
		return res
	}

	publishedAt, err := time.Parse(time.RFC3339, req.PublishedAt)
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid published_at: %w", err))
	}

	review, err := s.store.GetReviewByUUID(r.Context(), uid)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return core.Err(http.StatusNotFound, fmt.Errorf("review not found"))
		}
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get review: %w", err))
	}

	review.Name = req.Name
	review.Content = req.Content
	review.FileID = req.FileID
	review.UserID = user.ID
	review.PublishedAt = publishedAt
	review.UpdatedAt = time.Now()

	err = s.store.UpdateReview(r.Context(), review)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to update review: %w", err))
	}

	updated, err := s.store.GetReviewByUUID(r.Context(), uid)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to load updated review: %w", err))
	}

	return core.JSON(http.StatusOK, updated)
}

func (s *Service) deleteReview(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	uid, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid review uuid: %w", err))
	}

	err = s.store.DeleteReviewByUUID(r.Context(), uid)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return core.Err(http.StatusNotFound, fmt.Errorf("review not found"))
		}
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to delete review: %w", err))
	}

	return core.JSON(http.StatusNoContent, nil)
}
