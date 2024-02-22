package httpmux

import "net/http"

type Mux struct {
	serveMux    *http.ServeMux
	NotFound    http.Handler
	middlewares []func(http.Handler) http.Handler
}

func New() *Mux {
	return &Mux{
		serveMux: &http.ServeMux{},
		NotFound: http.NotFoundHandler(),
	}
}

func (m *Mux) wrap(handler http.Handler) http.Handler {
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		handler = m.middlewares[i](handler)
	}

	return handler
}

func (m *Mux) Handle(pattern string, handler http.Handler) {
	m.serveMux.Handle(pattern, m.wrap(handler))
}

func (m *Mux) HandleFunc(pattern string, handler http.HandlerFunc) {
	m.Handle(pattern, m.wrap(handler))
}

func (m *Mux) Handler(r *http.Request) (h http.Handler, pattern string) {
	return m.serveMux.Handler(r)
}

func (m *Mux) Use(mw ...func(http.Handler) http.Handler) {
	m.middlewares = append(m.middlewares, mw...)
}

func (m *Mux) Group(fn func(*Mux)) {
	mux := *m
	fn(&mux)
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, pattern := m.serveMux.Handler(r)
	if pattern == "" {
		m.NotFound.ServeHTTP(w, r)
		return
	}

	m.serveMux.ServeHTTP(w, r)
}
