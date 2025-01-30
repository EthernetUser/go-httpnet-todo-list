package router

import (
	"net/http"
	"strings"
)

type Router interface {
	GetMux() *http.ServeMux
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request))
	Get(pattern string, handler func(w http.ResponseWriter, r *http.Request))
	Post(pattern string, handler func(w http.ResponseWriter, r *http.Request))
	Put(pattern string, handler func(w http.ResponseWriter, r *http.Request))
	Delete(pattern string, handler func(w http.ResponseWriter, r *http.Request))
}

type router struct {
	m *http.ServeMux
}

type Middleware func(next http.Handler) http.Handler

func New() Router {
	return &router{
		m: http.NewServeMux(),
	}
}

func (r *router) GetMux() *http.ServeMux {
	return r.m
}

func (r *router) Handle(pattern string, handler http.Handler) {
	r.m.Handle(pattern, handler)
}

func (r *router) HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	r.m.HandleFunc(pattern, handler)
}

func (r *router) Get(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	updatedPattern := "GET " + strings.TrimSpace(pattern)
	r.m.HandleFunc(updatedPattern, handler)
}

func (r *router) Post(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	updatedPattern := "POST " + strings.TrimSpace(pattern)
	r.m.HandleFunc(updatedPattern, handler)
}

func (r *router) Put(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	updatedPattern := "PUT " + strings.TrimSpace(pattern)
	r.m.HandleFunc(updatedPattern, handler)
}

func (r *router) Delete(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	updatedPattern := "DELETE " + strings.TrimSpace(pattern)
	r.m.HandleFunc(updatedPattern, handler)
}

func AddMiddlewares(middleware ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for _, m := range middleware {
			next = m(next)
		}
		return next
	}
}
