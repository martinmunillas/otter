package server

import "net/http"

type Middleware = func(next http.Handler) http.Handler

func (s *server) Use(middleware Middleware) *server {
	s.middlewares = append(s.middlewares, middleware)
	return s
}
