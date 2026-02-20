package core

import (
	"encoding/json/v2"
	"log/slog"
	"net/http"

	"github.com/zagvozdeen/ola/internal/logger"
	"github.com/zagvozdeen/ola/internal/store/models"
)

type HandlerFunc func(*http.Request, *models.User) Response

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

//type FlushData struct {
//	ctx  context.Context
//	data <-chan []byte
//}

var _ Response = (*ResponseError)(nil)
var _ Response = (*ResponseData)(nil)

//var _ Response = (*FlushData)(nil)

func Err(code int, err error) *ResponseError {
	return &ResponseError{code: code, err: err}
}

func JSON(code int, d any) *ResponseData {
	return &ResponseData{code: code, data: d}
}

//func Flush(ctx context.Context, data <-chan []byte) *FlushData {
//	return &FlushData{ctx: ctx, data: data}
//}

func (r *ResponseError) Response(w http.ResponseWriter, log *logger.Logger) int {
	log.Debug("Internal error", slog.Any("error", r.err), slog.Int("code", r.code))
	http.Error(w, r.err.Error(), r.code)
	return r.code
}

func (r *ResponseData) Response(w http.ResponseWriter, log *logger.Logger) int {
	w.WriteHeader(r.code)
	if r.code == http.StatusNoContent && r.data == nil {
		return r.code
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.MarshalWrite(w, r.data)
	if err != nil {
		log.Error("Failed to marshal response", err, slog.Int("code", r.code))
	}
	return r.code
}

//func (r *FlushData) Response(w http.ResponseWriter, log *logger.Logger) int {
//	flusher, ok := w.(http.Flusher)
//	if !ok {
//		return Err(http.StatusHTTPVersionNotSupported, errors.New("streaming not supported")).Response(w, log)
//	}
//
//	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
//	w.Header().Set("Cache-Control", "no-cache")
//	w.Header().Set("Connection", "keep-alive")
//
//	var b []byte
//	for {
//		select {
//		case <-r.ctx.Done():
//			return http.StatusGone
//		case b, ok = <-r.data:
//			if !ok {
//				return http.StatusOK
//			}
//			if _, err := w.Write(b); err != nil {
//				log.Error("Failed to write a piece of data", slog.Any("err", err))
//				return http.StatusInternalServerError
//			}
//			flusher.Flush()
//		}
//	}
//}
