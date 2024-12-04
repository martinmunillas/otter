package server

import (
	"fmt"
	"net/http"

	"github.com/martinmunillas/otter/server/tools"
)

type Handler = func(r *http.Request, t tools.Tools)

type Page struct {
	Path    string
	Handler Handler
}

func NewPage(path string, handler Handler) Page {
	return Page{
		Path:    path,
		Handler: handler,
	}
}

func (s *Server) HandlePages(pages ...Page) *Server {
	for _, page := range pages {
		s.mux.Handle(fmt.Sprintf("GET %s", page.Path), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			page.Handler(r, tools.Make(w, r))
		}))
	}
	return s
}
