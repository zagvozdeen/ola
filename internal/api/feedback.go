package api

import (
	"net/http"

	"github.com/zagvozdeen/ola/internal/api/core"
	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Service) getFeedback(r *http.Request, user *models.User) core.Response {
	panic("implement")
}

func (s *Service) createFeedback(r *http.Request, user *models.User) core.Response {
	panic("implement")
}

type createGuestFeedbackRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Content string `json:"content"`
}

func (s *Service) createGuestFeedback(r *http.Request) core.Response {
	panic("implement")
}
