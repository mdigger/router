package mux

import (
	"context"
	"net/http"
	"strings"

	"github.com/mdigger/router"
)

type ServeMux struct {
	Headers  map[string]string
	Redirect func(w http.ResponseWriter, r *http.Request, status int, url string)
	NotFound func(w http.ResponseWriter, r *http.Request)
	Error    func(w http.ResponseWriter, status int, err error)
	routers  map[string]*router.Paths
}

func (m *ServeMux) Handle(method, path string, handler http.Handler) {
	if method == "" || handler == nil {
		return
	}
	if m.routers == nil {
		// typically no more than 9 of HTTP methods
		m.routers = make(map[string]*router.Paths, 9)
	}
	method = strings.ToUpper(method)
	r := m.routers[method]
	if r == nil {
		r = new(router.Paths)
		m.routers[method] = r
	}
	if err := r.Add(path, handler); err != nil {
		panic(err) // the handler does not suit us for some reason
	}
}

func (m *ServeMux) HandleFunc(method, path string,
	handleFunc func(http.ResponseWriter, *http.Request)) {
	m.Handle(method, path, http.HandlerFunc(handleFunc))
}

func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	responseHeader := w.Header()
	if len(m.Headers) > 0 {
		for key, value := range m.Headers {
			responseHeader.Set(key, value)
		}
	}
	var path = r.URL.Path
	// ctx := log.WithField("path", path)
	if routers := m.routers[r.Method]; routers != nil {
		if handler, params := routers.Lookup(path); handler != nil {
			if len(params) > 0 {
				ctx := context.WithValue(r.Context(), keyParams, params)
				r = r.WithContext(ctx)
			}
			// name := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
			// ctx.WithField("name", name).Debug("found handler")
			handler.(http.Handler).ServeHTTP(w, r)
			return
		}
		// handler not found
		if strings.HasSuffix(path, "/") {
			path = strings.TrimSuffix(path, "/")
		} else {
			path += "/"
		}
		if handler, _ := routers.Lookup(path); handler != nil {
			status := http.StatusMovedPermanently
			if r.Method != http.MethodGet && r.Method != http.MethodHead {
				status = http.StatusPermanentRedirect
			}
			// name := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
			// ctx.WithField("name", name).Debug("redirect to handler")
			if m.Redirect != nil {
				m.Redirect(w, r, status, path)
			} else {
				http.Redirect(w, r, path, status)
			}
			return
		}
	}
	// handler for request method not found
	var methods = make([]string, 0, len(m.routers))
	for method, handlers := range m.routers {
		if handler, _ := handlers.Lookup(path); handler != nil {
			methods = append(methods, method)
		}
	}
	if len(methods) > 0 {
		responseHeader.Set("Allow", strings.Join(methods, ", "))
		status := http.StatusMethodNotAllowed
		if m.Error != nil {
			m.Error(w, status, nil)
		} else {
			http.Error(w, http.StatusText(status), status)
		}
		return
	}
	// not found
	if m.NotFound != nil {
		m.NotFound(w, r)
	} else {
		http.NotFound(w, r)
	}
}

type contextKey byte // context key type
const (
	keyParams contextKey = iota
)

func PathParams(r *http.Request) router.Params {
	if params, ok := r.Context().Value(keyParams).(router.Params); ok {
		return params
	}
	return nil
}

var Default = new(ServeMux)

func Handle(method, path string, handler http.Handler) {
	Default.Handle(method, path, handler)
}

func HandleFunc(method, path string,
	handleFunc func(http.ResponseWriter, *http.Request)) {
	Default.HandleFunc(method, path, handleFunc)
}
