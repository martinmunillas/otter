package server

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"path/filepath"
)

func (s *Server) ServeStatic(path string) *Server {
	staticDir, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	s.mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(staticDir))))

	slog.Info(fmt.Sprintf("Serving static files from %s", staticDir))
	return s
}
