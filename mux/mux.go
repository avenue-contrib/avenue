package mux

import (
	"net/http"
	"path"

	"github.com/avenue-contrib/avenue/context"

	"github.com/julienschmidt/httprouter"
)

type Mux struct {
	Handlers []context.Handler
	router   *httprouter.Router
	prefix   string
}

func New(prefix string) Mux {
	return Mux{
		Handlers: make([]context.Handler, 0),
		router:   httprouter.New(),
		prefix:   prefix,
	}
}

func (m *Mux) Router() *httprouter.Router {
	return m.router
}

func (m *Mux) Handle(method, p string, handlers []context.Handler) {
	p = path.Join(m.prefix, p)
	m.router.Handle(method, p, func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		// create a new context
		ctx := context.New(*req, res)

		// store the url parameters -- be careful not to overwrite with .Set()
		ctx.Parameters = params

		// give the context a handler chain
		// TODO: determine allocation hit of this array-per-context
		ctx.Handlers = append(m.Handlers, handlers...)

		// start things off
		ctx.Next()
	})
}

func (m *Mux) POST(path string, handlers ...context.Handler) {
	m.Handle("POST", path, handlers)
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (m *Mux) GET(path string, handlers ...context.Handler) {
	m.Handle("GET", path, handlers)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (m *Mux) DELETE(path string, handlers ...context.Handler) {
	m.Handle("DELETE", path, handlers)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (m *Mux) PATCH(path string, handlers ...context.Handler) {
	m.Handle("PATCH", path, handlers)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (m *Mux) PUT(path string, handlers ...context.Handler) {
	m.Handle("PUT", path, handlers)
}
