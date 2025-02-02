package router

import (
	"net/http"
	"strings"
)

type Router interface {
	GetMux() *http.ServeMux
	Handle(pattern string, handler http.Handler)
	HandleFunc(
		pattern string,
		handler func(w http.ResponseWriter, r *http.Request),
	)
	Get(pattern string, handler func(w http.ResponseWriter, r *http.Request))
	Post(pattern string, handler func(w http.ResponseWriter, r *http.Request))
	Put(pattern string, handler func(w http.ResponseWriter, r *http.Request))
	Delete(pattern string, handler func(w http.ResponseWriter, r *http.Request))
}

const (
	EmptySpaceString = " "
)

type router struct {
	m *http.ServeMux
}

type Middleware func(next http.Handler) http.Handler

func New() Router {
	return &router{
		m: http.NewServeMux(),
	}
}

func updatePatternWithMethod(method string, pattern string) string {
	trimedPattern := strings.TrimSpace(pattern)
	return strings.Join([]string{method, trimedPattern}, EmptySpaceString)
}

func (r *router) GetMux() *http.ServeMux {
	return r.m
}

func (r *router) Handle(pattern string, handler http.Handler) {
	r.m.Handle(pattern, handler)
}

func (r *router) HandleFunc(
	pattern string,
	handler func(w http.ResponseWriter, r *http.Request),
) {
	r.m.HandleFunc(pattern, handler)
}

func (r *router) Get(
	pattern string,
	handler func(w http.ResponseWriter, r *http.Request),
) {
	updatedPattern := updatePatternWithMethod(http.MethodGet, pattern)
	r.m.HandleFunc(updatedPattern, handler)
}

func (r *router) Post(
	pattern string,
	handler func(w http.ResponseWriter, r *http.Request),
) {
	updatedPattern := updatePatternWithMethod(http.MethodPost, pattern)
	r.m.HandleFunc(updatedPattern, handler)
}

func (r *router) Put(
	pattern string,
	handler func(w http.ResponseWriter, r *http.Request),
) {
	updatedPattern := updatePatternWithMethod(http.MethodPut, pattern)
	r.m.HandleFunc(updatedPattern, handler)
}

func (r *router) Delete(
	pattern string,
	handler func(w http.ResponseWriter, r *http.Request),
) {
	updatedPattern := updatePatternWithMethod(http.MethodDelete, pattern)
	r.m.HandleFunc(updatedPattern, handler)
}

func CreateMiddlewaresWrapper(middleware ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for _, m := range middleware {
			next = m(next)
		}
		return next
	}
}
