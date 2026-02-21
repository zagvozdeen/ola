package api

import (
	"html/template"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
)

type PageData struct {
	InsertRootDiv bool
	Head          template.HTML
}

func (s *Service) index(w http.ResponseWriter, r *http.Request) {
	if !s.cfg.IsProduction && slices.Contains([]string{".ico", ".png", ".jpg", ".jpeg", ".gif", ".svg"}, filepath.Ext(r.URL.Path)) {
		http.ServeFile(w, r, "web/public"+r.URL.Path)
		return
	}
	if !s.cfg.IsProduction && strings.HasPrefix(r.URL.Path, "/files/") {
		http.ServeFile(w, r, ".data"+r.URL.Path)
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
	err = tmlp.Execute(w, data)
	if err != nil {
		s.log.Error("Failed to execute template", err)
		return
	}
}
