package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/martinmunillas/otter"
	"github.com/martinmunillas/otter/i18n"
	"github.com/martinmunillas/otter/response/send"
)

type server struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

type SendOk struct {
	HTML func(component templ.Component)
	JSON func(content any)
}
type SendUnauthorized struct {
	HTML func(component templ.Component)
	JSON func(message string)
}
type SendForbidden struct {
	HTML func(component templ.Component)
	JSON func(message string)
}
type SendNotFound struct {
	HTML func(component templ.Component)
	JSON func(message string)
}
type SendBadRequest struct {
	HTML func(component templ.Component)
	JSON func(message string)
}
type SendInternalError struct {
	HTML func(err error, component templ.Component)
	JSON func(err error)
}

type Send struct {
	Ok            SendOk
	Unauthorized  SendUnauthorized
	Forbidden     SendForbidden
	NotFound      SendNotFound
	BadRequest    SendBadRequest
	InternalError SendInternalError
}

type Redirect struct {
	Server func(path string, status int)
	HX     func(path string)
}

type Tools struct {
	T             func(key string, replacements ...i18n.Replacements) templ.Component
	RawT          func(key string, replacements ...i18n.Replacements) templ.Component
	Translation   func(key string) string
	ErrorT        func(key string) error
	DateTime      func(t time.Time, style i18n.DateStyle) string
	Send          Send
	Redirect      Redirect
	SetRawCookies func(rawCookies string)
	SetCookie     func(cookie http.Cookie)
	SetToast      func(toast otter.Toast)
}

type Handler = func(r *http.Request, t Tools)

type Endpoint struct {
	Method  string
	Path    string
	Handler Handler
}

func makeTools(w http.ResponseWriter, r *http.Request) Tools {
	ctx := r.Context()
	return Tools{
		T: func(key string, replacements ...i18n.Replacements) templ.Component {
			return i18n.T(ctx, key, replacements...)
		},
		RawT: func(key string, replacements ...i18n.Replacements) templ.Component {
			return i18n.RawT(ctx, key, replacements...)
		},
		Translation: func(key string) string { return i18n.Translation(ctx, key) },
		ErrorT:      func(key string) error { return i18n.ErrorT(ctx, key) },
		DateTime: func(t time.Time, style i18n.DateStyle) string {
			return i18n.DateTime(ctx, t, style)
		},
		Redirect: Redirect{
			Server: func(path string, status int) {
				http.Redirect(w, r, path, status)
			},
			HX: func(path string) {
				w.Header().Set("HX-Redirect", path)
			},
		},
		SetRawCookies: func(raw string) {
			w.Header().Set("Set-Cookie", raw)
		},
		SetCookie: func(cookie http.Cookie) {
			http.SetCookie(w, &cookie)
		},
		SetToast: func(toast otter.Toast) {
			eventMap := map[string]otter.Toast{}
			eventMap["makeToast"] = toast
			jsonData, err := json.Marshal(eventMap)
			if err != nil {
				return
			}
			w.Header().Set("HX-Trigger", string(jsonData))
		},
		Send: Send{
			Ok: SendOk{
				HTML: func(component templ.Component) {
					send.Html.Ok(w, ctx, component)
				},
				JSON: func(content any) {
					send.Json.Ok(w, content)
				},
			},
			Unauthorized: SendUnauthorized{
				HTML: func(component templ.Component) {
					send.Html.Unauthorized(w, ctx, component)
				},
				JSON: func(message string) {
					send.Json.Unauthorized(w, message)
				},
			},
			Forbidden: SendForbidden{
				HTML: func(component templ.Component) {
					send.Html.Forbidden(w, ctx, component)
				},
				JSON: func(message string) {
					send.Json.Forbidden(w, message)
				},
			},
			NotFound: SendNotFound{
				HTML: func(component templ.Component) {
					send.Html.NotFound(w, ctx, component)
				},
				JSON: func(message string) {
					send.Json.NotFound(w, message)
				},
			},
			BadRequest: SendBadRequest{
				HTML: func(component templ.Component) {
					send.Html.BadRequest(w, ctx, component)
				},
				JSON: func(message string) {
					send.Json.BadRequest(w, message)
				},
			},
			InternalError: SendInternalError{
				HTML: func(err error, component templ.Component) {
					send.Html.InternalError(w, ctx, err, component)
				},
				JSON: func(err error) {
					send.Json.InternalError(w, err)
				},
			},
		},
	}
}

func NewServer(endpoints []Endpoint) *server {
	mux := http.NewServeMux()
	for _, endpoint := range endpoints {
		mux.Handle(fmt.Sprintf("%s %s", endpoint.Method, endpoint.Path), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			endpoint.Handler(r, makeTools(w, r))
		}))
	}
	return &server{
		mux: mux,
	}
}

func (s server) Listen(port int64) {
	slog.Info(fmt.Sprintf("Server listening on port %d", port))
	handler := i18n.Middleware(s.mux)
	for _, middleware := range s.middlewares {
		handler = middleware(handler)
	}
	err := http.ListenAndServe(PortString(port), handler)
	if err != nil {
		slog.Error(err.Error())
	}
}
