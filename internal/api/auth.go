package api

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/zagvozdeen/ola/internal/api/core"
	"github.com/zagvozdeen/ola/internal/store/models"
	"golang.org/x/crypto/bcrypt"
)

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
}

func (s *Service) login(r *http.Request) core.Response {
	req := &authRequest{}
	err := json.UnmarshalRead(r.Body, req)
	if err != nil {
		s.log.Warn("Failed to decode auth request", slog.Any("err", err))
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	user, err := s.store.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warn("Invalid credentials", slog.String("username", req.Username))
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}
		s.log.Error("Failed to load user", err, slog.String("username", req.Username))
		http.Error(w, "failed to authenticate", http.StatusInternalServerError)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(req.Password))
	if err != nil {
		s.log.Warn("Invalid credentials", slog.String("username", req.Username))
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        strconv.Itoa(user.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 365)),
	})
	token, err := t.SignedString([]byte(s.cfg.AppSecret))
	if err != nil {
		s.log.Error("Failed to sign auth token", err)
		http.Error(w, "failed to create token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.MarshalWrite(w, authResponse{Token: token})
	if err != nil {
		s.log.Error("Failed to write response", err)
	}
}

type registerRequest struct {
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type registerResponse struct {
	UUID      string    `json:"uuid"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Service) register(r *http.Request) core.Response {
	req := &registerRequest{}
	err := json.UnmarshalRead(r.Body, req)
	if err != nil {
		return core.Err(http.StatusBadRequest, err)
	}
	err = s.validate.StructCtx(r.Context(), req)
	if err != nil {
		return core.Err(http.StatusBadRequest, err)
	}
	return core.JSON(http.StatusCreated, registerResponse{
		UUID:      "",
		Email:     "",
		CreatedAt: time.Time{},
	})
}

func (s *Service) guest(fn core.GuestHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(r).Response(w, s.log)
	}
}

func (s *Service) auth(fn core.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, res := s.checkAuth(r, r.Header.Get("Authorization"))
		if res == nil {
			res = fn(r, user)
		}
		res.Response(w, s.log)
	}
}

func (s *Service) checkAuth(r *http.Request, token string) (*models.User, core.Response) {
	switch {
	case strings.HasPrefix(token, "tma "):
		return s.authTMA(r, token)
	case strings.HasPrefix(token, "Bearer "):
		return s.authBearer(r, token)
	default:
		return nil, core.Err(http.StatusUnauthorized, fmt.Errorf("missing authorization header"))
	}
}

func (s *Service) authTMA(r *http.Request, token string) (*models.User, core.Response) {
	token = strings.TrimPrefix(token, "tma ")
	values, err := url.ParseQuery(token)
	if err != nil {
		return nil, core.Err(http.StatusUnauthorized, fmt.Errorf("failed to parse tma token: %w", err))
	}
	u, ok := bot.ValidateWebappRequest(values, s.cfg.TelegramBotToken)
	if !ok {
		return nil, core.Err(http.StatusUnauthorized, fmt.Errorf("invalid tma token"))
	}
	var user *models.User
	user, err = s.store.GetUserByTID(r.Context(), u.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.Err(http.StatusUnauthorized, fmt.Errorf("tma user not found: %w", err))
		}
		return nil, core.Err(http.StatusInternalServerError, fmt.Errorf("failed to load user: %w", err))
	}
	return user, nil
}

func (s *Service) authBearer(r *http.Request, token string) (*models.User, core.Response) {
	token = strings.TrimPrefix(token, "Bearer ")
	var claims jwt.RegisteredClaims
	t, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (any, error) {
		return []byte(s.cfg.AppSecret), nil
	})
	if err != nil {
		return nil, core.Err(http.StatusUnauthorized, fmt.Errorf("failed to parse token: %w", err))
	}
	if !t.Valid {
		return nil, core.Err(http.StatusUnauthorized, fmt.Errorf("invalid token"))
	}
	id, err := strconv.Atoi(claims.ID)
	if err != nil {
		return nil, core.Err(http.StatusUnauthorized, fmt.Errorf("invalid token: %w, id=%s", err, claims.ID))
	}
	var user *models.User
	user, err = s.store.GetUserByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.Err(http.StatusUnauthorized, fmt.Errorf("user not found: %w", err))
		}
		return nil, core.Err(http.StatusInternalServerError, fmt.Errorf("failed to load user: %w", err))
	}
	return user, nil
}
