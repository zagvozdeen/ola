package api

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/zagvozdeen/ola/internal/config"
	"github.com/zagvozdeen/ola/internal/logger"
	"github.com/zagvozdeen/ola/internal/store"
)

type Service struct {
	cfg          *config.Config
	log          *logger.Logger
	store        store.Storage
	processingTS sync.Map
	bot          *bot.Bot
	botStarted   chan struct{}
}

func New(cfg *config.Config, log *logger.Logger, store store.Storage) *Service {
	return &Service{
		cfg:          cfg,
		log:          log,
		store:        store,
		processingTS: sync.Map{},
		botStarted:   make(chan struct{}, 1),
	}
}

func (s *Service) Run(ctx context.Context) {
	server := &http.Server{
		Addr:     net.JoinHostPort(s.cfg.APIHost, s.cfg.APIPort),
		Handler:  s.getRoutes(),
		ErrorLog: s.log.GetLog(),
	}

	errCh := make(chan error, 1)
	wg := &sync.WaitGroup{}
	wg.Go(func() {
		errCh <- server.ListenAndServe()
		close(errCh)
	})
	s.log.Infof("Server started on %s", net.JoinHostPort(s.cfg.APIHost, s.cfg.APIPort))

	select {
	case <-ctx.Done():
		s.log.Info("Context canceled")
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			s.log.Info("Server has been closed")
			return
		}
		s.log.Error("Failed to listen and serve server", err)
		return
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	err := server.Shutdown(shutdownCtx)
	cancel()
	if err != nil {
		s.log.Error("Failed to shutdown server", err)
	}
	wg.Wait()
	s.log.Info("Service has been stopped")
}

func (s *Service) getRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	if !s.cfg.IsProduction {
		mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "dev.html")
		})
		mux.Handle("GET /node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("node_modules"))))
	}

	mux.HandleFunc("POST /api/auth", s.login)
	//mux.HandleFunc("GET /api/test-sessions", s.auth(s.getTestSessions))
	//mux.HandleFunc("GET /api/test-sessions/{uuid}", s.auth(s.getTestSession))
	//mux.HandleFunc("POST /api/test-sessions", s.auth(s.createTestSession))
	//mux.HandleFunc("PATCH /api/user-answers/{uuid}", s.auth(s.updateUserAnswer))
	//mux.HandleFunc("GET /api/leaderboard", s.auth(s.getLeaderboard))
	//mux.HandleFunc("GET /api/cards", s.auth(s.getCards))
	//mux.HandleFunc("GET /api/courses", s.auth(s.getCourses))
	//mux.HandleFunc("GET /api/modules", s.auth(s.getModules))
	//mux.HandleFunc("GET /api/changes", s.auth(s.getChanges))

	return mux
}
