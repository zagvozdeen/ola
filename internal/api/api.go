package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"

	"github.com/go-playground/mold/v4"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
	"github.com/go-telegram/bot"
	"github.com/zagvozdeen/ola/internal/config"
	"github.com/zagvozdeen/ola/internal/event_bus"
	"github.com/zagvozdeen/ola/internal/logger"
	"github.com/zagvozdeen/ola/internal/seeder"
	"github.com/zagvozdeen/ola/internal/store"
	"github.com/zagvozdeen/ola/internal/worker_pool"
)

type Service struct {
	cfg        *config.Config
	log        *logger.Logger
	store      *store.Store
	viteProxy  *httputil.ReverseProxy
	validate   *validator.Validate
	conform    *mold.Transformer
	workerPool *worker_pool.WorkerPool
	eventBus   *event_bus.EventBus
	bot        *bot.Bot
}

func New(cfg *config.Config, log *logger.Logger, store *store.Store) *Service {
	workerPool := worker_pool.New(log, 4, 100)
	return &Service{
		cfg:        cfg,
		log:        log,
		store:      store,
		viteProxy:  newViteProxy(log),
		validate:   newValidator(log),
		conform:    modifiers.New(),
		workerPool: workerPool,
		eventBus:   event_bus.New(workerPool),
	}
}

func (s *Service) Run(ctx context.Context) {
	addr := fmt.Sprintf("%s:%d", s.cfg.App.Host, s.cfg.App.Port)
	server := &http.Server{
		Addr:     addr,
		Handler:  s.getRoutes(),
		ErrorLog: s.log.GetLog(),
	}

	err := seeder.New(s.cfg, s.log, s.store).Run(ctx)
	if err != nil {
		s.log.Error("Failed to run seeder", err)
		return
	}

	s.registerListeners()

	errCh := make(chan error, 2)
	wg := &sync.WaitGroup{}
	wg.Go(func() {
		errCh <- server.ListenAndServe()
	})
	wg.Go(func() {
		errCh <- s.startBot(ctx)
	})
	wg.Go(func() {
		s.workerPool.Run(ctx)
	})
	s.log.Infof("Server started on %s", addr)

	select {
	case <-ctx.Done():
		s.log.Info("Context canceled")
	case err = <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			s.log.Info("Server has been closed")
			return
		}
		s.log.Error("Failed to listen and serve server", err)
		return
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	err = server.Shutdown(shutdownCtx)
	cancel()
	if err != nil {
		s.log.Error("Failed to shutdown server", err)
	}
	wg.Wait()
	close(errCh)
	s.log.Info("Service has been stopped")
}

func (s *Service) getRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.index)

	mux.HandleFunc("POST /api/auth/login", s.guest(s.login))
	mux.HandleFunc("POST /api/auth/register", s.guest(s.register))

	mux.HandleFunc("POST /api/guest/feedback", s.guest(s.createGuestFeedback))
	mux.HandleFunc("POST /api/guest/orders", s.guest(s.createGuestOrder))

	mux.HandleFunc("GET /api/me", s.auth(s.getMe))
	mux.HandleFunc("GET /api/products", s.auth(s.getProducts))
	mux.HandleFunc("POST /api/products", s.auth(s.createProduct))
	mux.HandleFunc("GET /api/products/{uuid}", s.auth(s.getProduct))
	mux.HandleFunc("PATCH /api/products/{uuid}", s.auth(s.updateProduct))
	mux.HandleFunc("DELETE /api/products/{uuid}", s.auth(s.deleteProduct))
	mux.HandleFunc("POST /api/files", s.auth(s.UploadFile))
	mux.HandleFunc("GET /api/categories", s.auth(s.getCategories))
	mux.HandleFunc("GET /api/categories/{uuid}", s.auth(s.getCategory))
	mux.HandleFunc("POST /api/categories", s.auth(s.createCategory))
	mux.HandleFunc("PATCH /api/categories/{uuid}", s.auth(s.updateCategory))
	mux.HandleFunc("DELETE /api/categories/{uuid}", s.auth(s.deleteCategory))
	mux.HandleFunc("GET /api/feedback", s.auth(s.getFeedback))
	mux.HandleFunc("GET /api/feedback/{uuid}", s.auth(s.getFeedbackByUUID))
	mux.HandleFunc("PATCH /api/feedback/{uuid}/status", s.auth(s.updateFeedbackStatus))
	mux.HandleFunc("POST /api/feedback", s.auth(s.createFeedback))
	mux.HandleFunc("GET /api/reviews", s.auth(s.getReviews))
	mux.HandleFunc("GET /api/reviews/{uuid}", s.auth(s.getReview))
	mux.HandleFunc("POST /api/reviews", s.auth(s.createReview))
	mux.HandleFunc("PATCH /api/reviews/{uuid}", s.auth(s.updateReview))
	mux.HandleFunc("DELETE /api/reviews/{uuid}", s.auth(s.deleteReview))
	mux.HandleFunc("GET /api/orders", s.auth(s.getOrders))
	mux.HandleFunc("GET /api/orders/{uuid}", s.auth(s.getOrder))
	mux.HandleFunc("PATCH /api/orders/{uuid}/status", s.auth(s.updateOrderStatus))
	mux.HandleFunc("POST /api/orders", s.auth(s.createOrder))
	mux.HandleFunc("POST /api/orders/from-cart", s.auth(s.createOrderFromCart))
	mux.HandleFunc("GET /api/cart", s.auth(s.getCart))
	mux.HandleFunc("POST /api/cart/items", s.auth(s.upsertCartItem))
	mux.HandleFunc("DELETE /api/cart/items/{product_uuid}", s.auth(s.deleteCartItem))
	mux.HandleFunc("GET /api/users", s.auth(s.getUsers))
	mux.HandleFunc("GET /api/users/{uuid}", s.auth(s.getUser))
	mux.HandleFunc("PATCH /api/users/{uuid}/role", s.auth(s.updateUserRole))

	return mux
}
