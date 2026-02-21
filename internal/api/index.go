package api

import (
	"html/template"
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
	Categories    []models.Category
	Reviews       []models.Review
}

func (s *Service) index(w http.ResponseWriter, r *http.Request) {
	if !s.cfg.IsProduction && strings.HasPrefix(r.URL.Path, "/files/") {
		http.ServeFile(w, r, ".data"+r.URL.Path)
		return
	}
	if !s.cfg.IsProduction && slices.Contains([]string{".ico", ".png", ".jpg", ".jpeg", ".gif", ".svg"}, filepath.Ext(r.URL.Path)) {
		http.ServeFile(w, r, "web/public"+r.URL.Path)
		return
	}
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
		s.log.Error("Failed to execute template", err)
		return
	}
}
