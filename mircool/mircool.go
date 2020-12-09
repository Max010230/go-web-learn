package mircool

import (
	"net/http"
)

type HandleFunc func(c *Context)

type Engine struct {
	router *Router
}

func NewServer() *Engine {
	return &Engine{router: NewRouter()}
}

func (e *Engine) addRoute(method, path string, handler HandleFunc) {
	e.router.addRouter(method, path, handler)
}

func (e *Engine) GET(path string, handler HandleFunc) {
	e.addRoute("GET", path, handler)
}

func (e *Engine) POST(path string, handler HandleFunc) {
	e.addRoute("POST", path, handler)
}

func (e *Engine) PUT(path string, handler HandleFunc) {
	e.addRoute("PUT", path, handler)
}

func (e *Engine) DELETE(path string, handler HandleFunc) {
	e.addRoute("DELETE", path, handler)
}

func (e *Engine) PATCH(path string, handler HandleFunc) {
	e.addRoute("PATCH", path, handler)
}

func (e *Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	context := newContext(resp, req)
	e.router.handle(context)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}
