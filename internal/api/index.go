package api

import "net/http"

func (s *Service) index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/landing/index.html")
}
