package api

import (
	"html/template"
	"net/http"
	"path/filepath"
	"slices"
)

type PageData struct {
	InsertRootDiv bool
	Head          template.HTML
}

func (s *Service) index(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r.URL.Path)
	if !s.cfg.IsProduction && slices.Contains([]string{".ico", ".png", ".jpg", ".jpeg", ".gif", ".svg"}, filepath.Ext(r.URL.Path)) {
		http.ServeFile(w, r, "web/public"+r.URL.Path)
		return
	}
	// if r.URL.Path == "/favicon.ico" {
	// 	http.ServeFile(w, r, "web/public/favicon.ico")
	// 	return
	// }

	// funcMap := template.FuncMap{
	// 	"image": func(name string) (string, error) {
	// 		return "http://localhost:5173/images/" + name, nil
	// 	},
	// }
	tmlp, err := template.ParseFiles("templates/index.html")
	if err != nil {
		s.log.Error("Failed to parse template", err)
		return
	}
	var head template.HTML
	var insertRootDiv bool
	// if r.URL.Path == "/" {
	head = `<script type="module" src="http://localhost:5173/@vite/client"></script>
    <script type="module" src="http://localhost:5173/landing/src/main.ts"></script>`
	// } else {
	// insertRootDiv = !s.cfg.IsProduction
	// head = `<script type="module" src="http://localhost:5174/@vite/client"></script>
	// <script type="module" src="http://localhost:5174/src/main.ts"></script>`
	// }
	data := PageData{
		InsertRootDiv: insertRootDiv,
		Head:          head,
	}
	err = tmlp.Execute(w, data)
	if err != nil {
		s.log.Error("Failed to execute template", err)
		return
	}
}
