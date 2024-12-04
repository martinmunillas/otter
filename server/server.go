package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/martinmunillas/otter/i18n"
)

type Server struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

func NewServer() *Server {
	return &Server{
		mux: http.NewServeMux(),
	}
}

func (s *Server) Listen(port int64) {
	isDevServer := os.Getenv("OTTER_DEV_SERVER") == "true"
	if isDevServer {
		slog.Info(fmt.Sprintf("Server listening on http://localhost:%d", port+1))
	} else {
		slog.Info(fmt.Sprintf("Server listening on port %d", port))
	}

	handler := i18n.Middleware(s.mux)
	for _, middleware := range s.middlewares {
		handler = middleware(handler)
	}
	err := http.ListenAndServe(PortString(port), handler)
	if err != nil {
		slog.Error(err.Error())
	}
}
