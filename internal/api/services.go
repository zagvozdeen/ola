package api

import (
	"net/http"

	"github.com/zagvozdeen/ola/internal/api/core"
	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Service) getServices(r *http.Request, user *models.User) core.Response {
	panic("implement")
}
