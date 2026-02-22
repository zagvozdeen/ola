package core

import (
	"encoding/json/v2"
	"log/slog"
	"net/http"

	"github.com/zagvozdeen/ola/internal/logger"
	"github.com/zagvozdeen/ola/internal/store/models"
)

type HandlerFunc func(*http.Request, *models.User) Response
type GuestHandlerFunc func(*http.Request) Response

type Response interface {
	Response(w http.ResponseWriter, log *logger.Logger) int
}

type ResponseError struct {
	code int
	err  error
}

type ResponseData struct {
	code int
	data any
}

var _ Response = (*ResponseError)(nil)
var _ Response = (*ResponseData)(nil)

func Err(code int, err error) *ResponseError {
	return &ResponseError{code: code, err: err}
}

func JSON(code int, d any) *ResponseData {
	return &ResponseData{code: code, data: d}
}

func (r *ResponseError) Response(w http.ResponseWriter, log *logger.Logger) int {
	log.Debug("Internal error", slog.Any("error", r.err), slog.Int("code", r.code))
	http.Error(w, r.err.Error(), r.code)
	return r.code
}

func (r *ResponseData) Response(w http.ResponseWriter, log *logger.Logger) int {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.code)
	if r.code == http.StatusNoContent && r.data == nil {
		return r.code
	}
	err := json.MarshalWrite(w, r.data)
	if err != nil {
		log.Error("Failed to marshal response", err, slog.Int("code", r.code))
	}
	return r.code
}
