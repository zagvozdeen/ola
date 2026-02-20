package api

import (
	"fmt"
	"html/template"
	"net/http"
)

type PageData struct {
	InsertRootDiv bool
	Head          template.HTML
}

func (s *Service) index(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	tmlp, err := template.ParseFiles("web/landing/index.html")
	if err != nil {
		s.log.Error("error 1", err)
		return
	}
	var head template.HTML
	var insertRootDiv bool
	if r.URL.Path == "/" {
		head = `<script type="module" src="http://localhost:5173/@vite/client"></script>
    <link rel="stylesheet" href="http://localhost:5173/index.css">
    <script type="module" src="http://localhost:5173/main.ts"></script>`
	} else {
		insertRootDiv = !s.cfg.IsProduction
		head = `<script type="module" src="http://localhost:5174/@vite/client"></script>
    <script type="module" src="http://localhost:5174/src/main.ts"></script>`
	}
	data := PageData{
		InsertRootDiv: insertRootDiv,
		Head:          head,
	}
	err = tmlp.Execute(w, data)
	if err != nil {
		s.log.Error("error 2", err)
		return
	}
}
