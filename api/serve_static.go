package api

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"path/filepath"
)

func ServeStatic(path string, mux *http.ServeMux) {
	staticDir, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(staticDir))))

	slog.Info(fmt.Sprintf("Serving static files from %s", staticDir))
}
