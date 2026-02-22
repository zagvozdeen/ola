package api

import (
	"encoding/json/v2"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/api/core"
	"github.com/zagvozdeen/ola/internal/store/models"
)

type createFeedbackRequest struct {
	Name    string `json:"name" mold:"trim" validate:"required,max=255"`
	Phone   string `json:"phone" mold:"trim" validate:"required,max=255,ru_phone"`
	Content string `json:"content" mold:"trim" validate:"required,max=3000"`
}

func (s *Service) getFeedback(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	feedback, err := s.store.GetAllFeedback(r.Context())
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get feedback: %w", err))
	}
	return core.JSON(http.StatusOK, feedback)
}

func (s *Service) createFeedback(r *http.Request, user *models.User) core.Response {
	req := &createFeedbackRequest{}
	err := json.UnmarshalRead(r.Body, req)
	if err != nil {
		return core.Err(http.StatusBadRequest, err)
	}
	err = s.conform.Struct(r.Context(), req)
	if err != nil {
		return core.Err(http.StatusBadRequest, err)
	}
	err = s.validate.StructCtx(r.Context(), req)
	if err != nil {
		return core.Err(http.StatusBadRequest, err)
	}

	uid, err := uuid.NewV7()
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to generate uuid v7: %w", err))
	}
	feedback := &models.Feedback{
		UUID:      uid,
		Name:      req.Name,
		Phone:     req.Phone,
		Content:   req.Content,
		UserID:    &user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = s.store.CreateFeedback(r.Context(), feedback)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to create feedback: %w", err))
	}

	return core.JSON(http.StatusCreated, feedback)
}

type createGuestFeedbackRequest struct {
	Name    string `json:"name" mold:"trim" validate:"required,max=255"`
	Phone   string `json:"phone" mold:"trim" validate:"required,max=255,ru_phone"`
	Content string `json:"content" mold:"trim" validate:"required,max=3000"`
}

func (s *Service) createGuestFeedback(r *http.Request) core.Response {
	req := &createGuestFeedbackRequest{}
	err := json.UnmarshalRead(r.Body, req)
	if err != nil {
		return core.Err(http.StatusBadRequest, err)
	}
	err = s.conform.Struct(r.Context(), req)
	if err != nil {
		return core.Err(http.StatusBadRequest, err)
	}
	err = s.validate.StructCtx(r.Context(), req)
	if err != nil {
		return core.Err(http.StatusBadRequest, err)
	}

	uid, err := uuid.NewV7()
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to generate uuid v7: %w", err))
	}
	feedback := &models.Feedback{
		UUID:      uid,
		Name:      req.Name,
		Phone:     req.Phone,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = s.store.CreateFeedback(r.Context(), feedback)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to create guest feedback: %w", err))
	}

	return core.JSON(http.StatusCreated, feedback)
}
