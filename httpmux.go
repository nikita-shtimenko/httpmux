package httpmux

import (
	"net/http"
)

// Mux is a http.Handler, which uses default http.ServeMux to dispatch requests to different handlers.
type Mux struct {
	serveMux    *http.ServeMux
	NotFound    http.Handler
	middlewares []func(http.Handler) http.Handler
}

// New returns a new initialized Mux instance.
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

// Handle registers the handler for the given pattern.
// If the given pattern conflicts, with one that is already registered, Handle
// panics.
func (m *Mux) Handle(pattern string, handler http.Handler) {
	m.serveMux.Handle(pattern, m.wrap(handler))
}

// HandleFunc registers the handler function for the given pattern.
// If the given pattern conflicts, with one that is already registered, HandleFunc
// panics.
func (m *Mux) HandleFunc(pattern string, handler http.HandlerFunc) {
	m.Handle(pattern, m.wrap(handler))
}

// Handler returns the handler to use for the given request,
// consulting r.Method, r.Host, and r.URL.Path. It always returns
// a non-nil handler. If the path is not in its canonical form, the
// handler will be an internally-generated handler that redirects
// to the canonical path. If the host contains a port, it is ignored
// when matching handlers.
//
// The path and host are used unchanged for CONNECT requests.
//
// Handler also returns the registered pattern that matches the
// request or, in the case of internally-generated redirects,
// the path that will match after following the redirect.
//
// If there is no registered handler that applies to the request,
// Handler returns a “page not found” handler and an empty pattern.
func (m *Mux) Handler(r *http.Request) (h http.Handler, pattern string) {
	return m.serveMux.Handler(r)
}

// Get registers the handler with HTTP GET method for the given pattern.
func (m *Mux) Get(pattern string, handler http.HandlerFunc) {
	m.Handle(http.MethodGet+" "+pattern, handler)
}

// Post registers the handler with HTTP POST method for the given pattern.
func (m *Mux) Post(pattern string, handler http.HandlerFunc) {
	m.Handle(http.MethodPost+" "+pattern, handler)
}

// Put registers the handler with HTTP PUT method for the given pattern.
func (m *Mux) Put(pattern string, handler http.HandlerFunc) {
	m.Handle(http.MethodPut+" "+pattern, handler)
}

// Delete registers the handler with HTTP DELETE method for the given pattern.
func (m *Mux) Delete(pattern string, handler http.HandlerFunc) {
	m.Handle(http.MethodDelete+" "+pattern, handler)
}

// Head registers the handler with HTTP HEAD method for the given pattern.
func (m *Mux) Head(pattern string, handler http.HandlerFunc) {
	m.Handle(http.MethodHead+" "+pattern, handler)
}

// Options registers the handler with HTTP OPTIONS method for the given pattern.
func (m *Mux) Options(pattern string, handler http.HandlerFunc) {
	m.Handle(http.MethodOptions+" "+pattern, handler)
}

// Use registers middleware with the Mux instance. Middleware must have the
// signature `func(http.Handler) http.Handler`.
func (m *Mux) Use(mw ...func(http.Handler) http.Handler) {
	m.middlewares = append(m.middlewares, mw...)
}

// Group is used to create 'groups' of routes in a Mux. Middleware registered
// inside the group will only be used on the routes in that group. See the
// example code at the start of the package documentation for how to use this
// feature.
func (m *Mux) Group(fn func(*Mux)) {
	mux := *m
	fn(&mux)
}

// ServeHTTP makes the router implement the http.Handler interface.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, pattern := m.serveMux.Handler(r)
	if pattern == "" {
		m.NotFound.ServeHTTP(w, r)
		return
	}

	m.serveMux.ServeHTTP(w, r)
}
