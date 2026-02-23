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
	Categories []models.Category
	Reviews    []models.Review
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
	if s.cfg.IsProduction {
		s.renderMainPage(w, r)
		return
	}
	if isWebSocketUpgrade(r) {
		s.viteProxy.ServeHTTP(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/admin") {
		r.URL.Path = "/spa/admin/"
		s.viteProxy.ServeHTTP(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/tma") {
		r.URL.Path = "/spa/tma/"
		s.viteProxy.ServeHTTP(w, r)
		return
	}
	vitePassThroughPrefixes := []string{
		"/@vite/",
		"/@id/",
		"/@fs/",
		//"/__vite_ping",
		"/shared/",
		//"/src/",         // if you ever reference /src directly
		//"/assets/",      // dev assets
		//"/landing/src/", // landing entry/modules
		"/spa/admin/src/",
		"/spa/tma/src/",
	}
	if strings.HasPrefix(r.URL.Path, "/files/") {
		http.ServeFile(w, r, ".data"+r.URL.Path)
		return
	}
	switch path.Ext(r.URL.Path) {
	case ".ico", ".png", ".jpg", ".jpeg", ".gif", ".svg", "webp":
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

	head := template.HTML(`<script type="module" src="http://localhost:5173/@vite/client"></script>
<script type="module" src="http://localhost:5173/landing/src/main.ts"></script>`)
	if s.cfg.IsProduction {
		head, err = s.renderViteHead(viteHeadParams{
			ManifestPath:         "public/.vite/manifest.json",
			EntryKey:             "landing/index.html",
			AssetsURLPrefix:      "/",
			IncludeCSS:           true,
			IncludeModulePreload: true,
		})
		if err != nil {
			s.log.Error("Failed to render vite head", err)
			return
		}
	}

	data := PageData{Head: head}
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
		s.log.Error("Failed to execute template", err, slog.String("proto", r.Proto), slog.String("url", r.URL.String()))
		return
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
