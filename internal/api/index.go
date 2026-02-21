package api

import (
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"slices"
	"strings"

	"github.com/zagvozdeen/ola/internal/store/models"
)

type PageData struct {
	InsertRootDiv bool
	Head          template.HTML
	Products      []models.Product
	Services      []models.Service
	Categories    []models.Category
	Reviews       []models.Review
}

func isWebSocketUpgrade(r *http.Request) bool {
	return strings.EqualFold(r.Header.Get("Upgrade"), "websocket") && strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade")
}

func (s *Service) index(w http.ResponseWriter, r *http.Request) {
	if s.cfg.IsProduction {
		s.renderMainPage(w, r)
		return
	}
	if isWebSocketUpgrade(r) {
		s.viteProxy.ServeHTTP(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/admin/") {
		r.URL.Path = "/spa/admin/"
		s.viteProxy.ServeHTTP(w, r)
	}
	if strings.HasPrefix(r.URL.Path, "/tma/") {
		r.URL.Path = "/spa/tma/"
		s.viteProxy.ServeHTTP(w, r)
	}
	vitePassThroughPrefixes := []string{
		"/@vite/",
		"/@id/",
		"/@fs/",
		"/__vite_ping",
		"/src/",         // if you ever reference /src directly
		"/assets/",      // dev assets
		"/landing/src/", // landing entry/modules
		"/spa/admin/src/",
		"/spa/tma/src/",
	}
	if !s.cfg.IsProduction && strings.HasPrefix(r.URL.Path, "/files/") {
		http.ServeFile(w, r, ".data"+r.URL.Path)
		return
	}
	if !s.cfg.IsProduction && slices.Contains([]string{".ico", ".png", ".jpg", ".jpeg", ".gif", ".svg"}, filepath.Ext(r.URL.Path)) {
		http.ServeFile(w, r, "web/public"+r.URL.Path)
		return
	}
	for _, pref := range vitePassThroughPrefixes {
		if strings.HasPrefix(r.URL.Path, pref) {
			s.viteProxy.ServeHTTP(w, r)
			return
		}
	}
	s.renderMainPage(w, r)
}

func (s *Service) renderMainPage(w http.ResponseWriter, r *http.Request) {
	tmlp, err := template.ParseFiles("templates/index.html")
	if err != nil {
		s.log.Error("Failed to parse template", err)
		return
	}
	data := PageData{
		InsertRootDiv: false,
		Head:          `<script type="module" src="http://localhost:5173/@vite/client"></script> <script type="module" src="http://localhost:5173/landing/src/main.ts"></script>`,
	}
	data.Products, err = s.store.GetAllProducts(r.Context())
	if err != nil {
		s.log.Error("Failed to get products", err)
		return
	}
	data.Services, err = s.store.GetAllServices(r.Context())
	if err != nil {
		s.log.Error("Failed to get services", err)
		return
	}
	data.Categories, err = s.store.GetAllCategories(r.Context())
	if err != nil {
		s.log.Error("Failed to get categories", err)
		return
	}
	data.Reviews, err = s.store.GetAllReviews(r.Context())
	if err != nil {
		s.log.Error("Failed to get all reviews", err)
		return
	}
	err = tmlp.Execute(w, data)
	if err != nil {
		s.log.Error("Failed to execute template", err, slog.String("proto", r.Proto), slog.String("url", r.URL.String()))
		return
	}
}
