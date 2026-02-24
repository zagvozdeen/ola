package api

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"

	"github.com/go-playground/mold/v4"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
	"github.com/zagvozdeen/ola/internal/config"
	"github.com/zagvozdeen/ola/internal/logger"
	"github.com/zagvozdeen/ola/internal/store"
)

type Service struct {
	cfg       *config.Config
	log       *logger.Logger
	store     *store.Store
	viteProxy *httputil.ReverseProxy
	validate  *validator.Validate
	conform   *mold.Transformer
}

func New(cfg *config.Config, log *logger.Logger, store *store.Store) *Service {
	return &Service{
		cfg:       cfg,
		log:       log,
		store:     store,
		viteProxy: newViteProxy(log),
		validate:  newValidator(log),
		conform:   modifiers.New(),
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

	mux.HandleFunc("GET /", s.index)

	mux.HandleFunc("POST /api/auth/login", s.guest(s.login))
	mux.HandleFunc("POST /api/auth/register", s.guest(s.register))

	mux.HandleFunc("POST /api/guest/feedback", s.guest(s.createGuestFeedback))
	mux.HandleFunc("POST /api/guest/orders", s.guest(s.createGuestOrder))

	mux.HandleFunc("GET /api/me", s.auth(s.getMe))                // for all
	mux.HandleFunc("GET /api/products", s.auth(s.getProducts))    // for all
	mux.HandleFunc("POST /api/products", s.auth(s.createProduct)) // for admin and moderator only
	mux.HandleFunc("GET /api/products/{uuid}", s.auth(s.getProduct))
	mux.HandleFunc("PATCH /api/products/{uuid}", s.auth(s.updateProduct))
	mux.HandleFunc("DELETE /api/products/{uuid}", s.auth(s.deleteProduct))
	mux.HandleFunc("POST /api/files", s.auth(s.UploadFile))        // for admin and moderator only
	mux.HandleFunc("GET /api/categories", s.auth(s.getCategories)) // for admin and moderator only
	mux.HandleFunc("GET /api/feedback", s.auth(s.getFeedback))     // for admin and moderator only
	mux.HandleFunc("POST /api/feedback", s.auth(s.createFeedback)) // for all
	mux.HandleFunc("GET /api/reviews", s.auth(s.getReviews))       // for admin and moderator only
	mux.HandleFunc("GET /api/orders", s.auth(s.getOrders))         // for admin and moderator only
	mux.HandleFunc("POST /api/orders", s.auth(s.createOrder))      // for all
	mux.HandleFunc("POST /api/orders/from-cart", s.auth(s.createOrderFromCart))
	mux.HandleFunc("GET /api/cart", s.auth(s.getCart))
	mux.HandleFunc("POST /api/cart/items", s.auth(s.upsertCartItem))
	mux.HandleFunc("DELETE /api/cart/items/{product_uuid}", s.auth(s.deleteCartItem))
	mux.HandleFunc("GET /api/users", s.auth(s.getUsers)) // for admin only

	return mux
}
