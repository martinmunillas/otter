package server

import "net/http"

type Middleware = func(next http.Handler) http.Handler

func (s *Server) Use(middleware Middleware) *Server {
	s.middlewares = append(s.middlewares, middleware)
	return s
}
