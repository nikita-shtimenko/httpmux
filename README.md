<div align="center">
httpmux - tiny wrapper on HTTP ServeMux for Go web applications
</div>

---

httpmux packs small set of features that you'll probably need:

- Create route **groups which use different middleware**.
- **Customizable handler** for `404 Not Found` response.
- Works with `http.Handler`, `http.HandlerFunc`, and standard Go middleware.
- Zero dependencies.
- Tiny and readable codebase (~60 lines of code).

---

### Installation

```
$ go get github.com/nikita-shtimenko/httpmux@latest
```

### Basic example

```go
mux := httpmux.New()

mux.NotFound = http.HandlerFunc(handlerNotFound)

// The Use() method can be used to register middleware.
// Middleware declared at the top level will used on all routes.
mux.Use(exampleMiddleware1)

// Default pattern registration
mux.HandleFunc("GET /api/v1/ping", exampleHandlerFunc1)

// You can create route 'groups'.
mux.Group(func(mux *httpmux.Mux) {
    // Middleware declared within in the group will only be used on the routes
    // in the group.
    mux.Use(exampleMiddleware2)

    mux.HandleFunc("GET /api/v1/status", exampleHandlerFunc2)

    // Groups can be nested.
    mux.Group(func(mux *httpmux.Mux) {
        mux.Use(exampleMiddleware3)

        mux.HandleFunc("GET /api/v1/users/{id}", exampleHandlerFunc3)
    })
})

func handlerNotFound(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(http.StatusText(http.StatusNotFound)))
}

func exampleHandlerFunc1(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("pong"))
}
```

### Notes

- It`s only a thin wrapper on default ServeMux, so the behaviour is default.
- Middleware must be declared _before_ a route in order to be used by that route. Any middleware declared after a route won't act on that route. For example:

```go
mux := httpmux.New()
mux.Use(middleware1)
mux.HandleFunc("GET /foo", ...) // This route will use middleware1 only.
mux.Use(middleware2)
mux.HandleFunc("POST /bar", ...) // This route will use both middleware1 and middleware2.
```

### Contributing

Bug fixes and documentation improvements are very welcome!
For feature additions or behavioral changes, please open an issue to discuss the change before submitting a PR.

### Thanks

Heavily inspired by Flow [alexedwards/flow](https://github.com/alexedwards/flow).
