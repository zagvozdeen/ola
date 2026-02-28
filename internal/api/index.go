package api

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/zagvozdeen/ola/internal/store/models"
)

type PageData struct {
	Head       template.HTML
	Products   []models.Product
	Services   []models.Product
	Categories []models.Category
	Title      string
	IsBlock    bool
}

type viteManifestEntry struct {
	File    string   `json:"file"`
	Imports []string `json:"imports"`
	CSS     []string `json:"css"`
}

type viteHeadParams struct {
	ManifestPath         string
	EntryKey             string
	AssetsURLPrefix      string
	IncludeCSS           bool
	IncludeModulePreload bool
}

func isWebSocketUpgrade(r *http.Request) bool {
	return strings.EqualFold(r.Header.Get("Upgrade"), "websocket") && strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade")
}

func (s *Service) index(w http.ResponseWriter, r *http.Request) {
	if s.cfg.App.IsProduction {
		s.renderMainPage(w, r)
		return
	}
	if isWebSocketUpgrade(r) {
		s.viteProxy.ServeHTTP(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/spa") {
		r.URL.Path = "/web/spa/"
		s.viteProxy.ServeHTTP(w, r)
		return
	}
	vitePassThroughPrefixes := []string{
		"/@vite/",
		"/@id/",
		"/@fs/",
		"/shared/",
		"/node_modules/",
		"/web/spa/src/",
		"/web/spa/src/",
	}
	if strings.HasPrefix(r.URL.Path, "/files/") {
		http.ServeFile(w, r, ".data"+r.URL.Path)
		return
	}
	switch path.Ext(r.URL.Path) {
	case ".ico", ".png", ".jpg", ".jpeg", ".gif", ".svg", "webp":
		http.ServeFile(w, r, "public"+r.URL.Path)
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
	isBlock, title, page, err := getTemplate(r)
	if err != nil {
		s.log.Warn("Get template error", slog.Any("error", err))
		http.NotFound(w, r)
		return
	}
	var templates *template.Template
	templates, err = s.getTemplates()
	if err != nil {
		s.log.Error("Get templates error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var head template.HTML
	if s.cfg.App.IsProduction {
		head, err = s.renderViteHead(viteHeadParams{
			ManifestPath:         "public/.vite/manifest.json",
			EntryKey:             "web/landing/index.html",
			AssetsURLPrefix:      "/",
			IncludeCSS:           true,
			IncludeModulePreload: true,
		})
		if err != nil {
			s.log.Error("Failed to render vite head", err)
			return
		}
	} else {
		head = `<script type="module" src="http://localhost:5173/@vite/client"></script> <script type="module" src="http://localhost:5173/web/landing/src/main.ts"></script>`
	}

	products, err := s.store.GetMainProducts(r.Context())
	if err != nil {
		s.log.Error("Failed to get main products", err)
		return
	}

	services, err := s.store.GetMainServices(r.Context())
	if err != nil {
		s.log.Error("Failed to get main services", err)
		return
	}
	categories, err := s.store.GetAllCategories(r.Context())
	if err != nil {
		s.log.Error("Failed to get categories", err)
		return
	}
	err = templates.ExecuteTemplate(w, page, PageData{
		Head:       head,
		Products:   products,
		Services:   services,
		Categories: categories,
		Title:      title,
		IsBlock:    isBlock,
	})
	if err != nil {
		s.log.Error("Failed to execute template", err, slog.String("path", r.URL.Path))
		return
	}
}

func (s *Service) getTemplates() (templates *template.Template, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.templates == nil {
		templates, err = template.ParseFiles("templates/index.html", "templates/delivery.html", "templates/privacy.html", "templates/templates.html")
		if err != nil {
			return nil, fmt.Errorf("failed to parse template: %w", err)
		}
	}
	if s.cfg.App.IsProduction {
		if s.templates == nil {
			s.templates = templates
		}
		return s.templates, nil
	}
	return templates, nil
}

func getTemplate(r *http.Request) (isBlock bool, title string, template string, err error) {
	switch r.URL.Path {
	case "/":
		return false, "OLA Studio", "index.html", nil
	case "/delivery", "/delivery/":
		return true, "Доставка, оплата и возврат | OLA Studio", "delivery.html", nil
	case "/catalog", "/catalog/":
		return true, "Каталог товаров и услуг | OLA Studio", "catalog.html", nil
	case "/privacy", "/privacy/":
		return true, "Политика конфиденциальности | OLA studio", "privacy.html", nil
	default:
		return false, "", "", fmt.Errorf("page not found: %s", r.URL.Path)
	}
}

func (s *Service) renderViteHead(params viteHeadParams) (template.HTML, error) {
	if params.ManifestPath == "" {
		return "", errors.New("vite manifest path is required")
	}
	if params.EntryKey == "" {
		return "", errors.New("vite manifest entry key is required")
	}

	prefix := params.AssetsURLPrefix
	if prefix == "" {
		prefix = "/"
	}
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	prefix = strings.TrimSuffix(prefix, "/")

	content, err := os.ReadFile(params.ManifestPath)
	if err != nil {
		return "", fmt.Errorf("read vite manifest: %w", err)
	}

	manifest := make(map[string]viteManifestEntry)
	if err = json.Unmarshal(content, &manifest); err != nil {
		return "", fmt.Errorf("parse vite manifest: %w", err)
	}

	entry, ok := manifest[params.EntryKey]
	if !ok {
		return "", fmt.Errorf("manifest entry %q not found", params.EntryKey)
	}
	if entry.File == "" {
		return "", fmt.Errorf("manifest entry %q has empty file", params.EntryKey)
	}

	formatAssetPath := func(asset string) string {
		asset = strings.TrimPrefix(asset, "/")
		return prefix + "/" + asset
	}

	links := make([]string, 0, len(entry.Imports))
	if params.IncludeModulePreload {
		visited := make(map[string]struct{})
		var collectImports func(importKeys []string) error
		collectImports = func(importKeys []string) error {
			for _, key := range importKeys {
				if _, exists := visited[key]; exists {
					continue
				}
				visited[key] = struct{}{}

				importEntry, exists := manifest[key]
				if !exists {
					return fmt.Errorf("manifest import %q not found", key)
				}
				if importEntry.File != "" {
					links = append(links, `<link rel="modulepreload" href=`+strconv.Quote(formatAssetPath(importEntry.File))+`>`)
				}
				if err = collectImports(importEntry.Imports); err != nil {
					return err
				}
			}
			return nil
		}

		if err = collectImports(entry.Imports); err != nil {
			return "", err
		}
	}

	if params.IncludeCSS {
		for _, css := range entry.CSS {
			links = append(links, `<link rel="stylesheet" href=`+strconv.Quote(formatAssetPath(css))+`>`)
		}
	}

	links = append(links, `<script type="module" src=`+strconv.Quote(formatAssetPath(entry.File))+`></script>`)
	return template.HTML(strings.Join(links, "\n")), nil
}
