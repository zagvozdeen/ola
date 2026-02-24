package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/api/core"
	"github.com/zagvozdeen/ola/internal/store/enums"
	"github.com/zagvozdeen/ola/internal/store/models"
)

type createFeedbackRequest struct {
	Name    string `json:"name" mold:"trim" validate:"required,max=255"`
	Phone   string `json:"phone" mold:"trim" validate:"required,max=255,ru_phone"`
	Content string `json:"content" mold:"trim" validate:"required,max=3000"`
	Type    string `json:"type" mold:"trim,lcase" validate:"required,oneof=manager_contact partnership_offer feedback_request"`
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

func (s *Service) getFeedbackByUUID(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	uid, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid feedback uuid: %w", err))
	}

	feedback, err := s.store.GetFeedbackByUUID(r.Context(), uid)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return core.Err(http.StatusNotFound, fmt.Errorf("feedback not found"))
		}
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get feedback: %w", err))
	}

	return core.JSON(http.StatusOK, feedback)
}

func (s *Service) createFeedback(r *http.Request, user *models.User) core.Response {
	req, res := core.Validate[createFeedbackRequest](r, s.conform, s.validate)
	if res != nil {
		return res
	}

	feedbackType, err := enums.NewFeedbackType(req.Type)
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid feedback type: %w", err))
	}

	src, ok := r.Context().Value("source").(enums.OrderSource)
	if !ok {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("source must be an enums.OrderSource"))
	}

	user.Phone = new(req.Phone)
	user.UpdatedAt = time.Now()
	err = s.store.UpdateUserPhone(r.Context(), user)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to update user phone: %w", err))
	}

	uid, err := uuid.NewV7()
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to generate uuid v7: %w", err))
	}
	feedback := &models.Feedback{
		UUID:      uid,
		Status:    enums.RequestStatusCreated,
		Source:    src,
		Type:      feedbackType,
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

	s.eventBus.FeedbackCreated.Publish(r.Context(), feedback)

	return core.JSON(http.StatusCreated, feedback)
}

type createGuestFeedbackRequest struct {
	Name    string `json:"name" mold:"trim" validate:"required,max=255"`
	Phone   string `json:"phone" mold:"trim" validate:"required,max=255,ru_phone"`
	Content string `json:"content" mold:"trim" validate:"required,max=3000"`
	Consent bool   `json:"consent" validate:"required"`
}

func (s *Service) createGuestFeedback(r *http.Request) core.Response {
	req, res := core.Validate[createGuestFeedbackRequest](r, s.conform, s.validate)
	if res != nil {
		return res
	}

	uid, err := uuid.NewV7()
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to generate uuid v7: %w", err))
	}
	feedback := &models.Feedback{
		UUID:      uid,
		Status:    enums.RequestStatusCreated,
		Source:    enums.OrderSourceLanding,
		Type:      enums.FeedbackTypeFeedbackRequest,
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

	s.eventBus.FeedbackCreated.Publish(r.Context(), feedback)

	return core.JSON(http.StatusCreated, feedback)
}

type updateFeedbackStatusRequest struct {
	Status string `json:"status" mold:"trim,lcase" validate:"required,oneof=created in_progress reviewed"`
}

func (s *Service) updateFeedbackStatus(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	uid, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid feedback uuid: %w", err))
	}

	req, res := core.Validate[updateFeedbackStatusRequest](r, s.conform, s.validate)
	if res != nil {
		return res
	}

	status, err := enums.NewRequestStatus(req.Status)
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid feedback status: %w", err))
	}

	feedback, err := s.store.GetFeedbackByUUID(r.Context(), uid)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return core.Err(http.StatusNotFound, fmt.Errorf("feedback not found"))
		}
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get feedback: %w", err))
	}

	err = s.store.UpdateFeedbackStatus(r.Context(), feedback.ID, status)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to update feedback status: %w", err))
	}

	updated, err := s.store.GetFeedbackByUUID(r.Context(), uid)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to load updated feedback: %w", err))
	}

	return core.JSON(http.StatusOK, updated)
}
