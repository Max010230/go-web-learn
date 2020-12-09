package mircool

import (
	"log"
	"net/http"
)

type Router struct {
	handlers map[string]HandleFunc
}

func NewRouter() *Router {
	return &Router{handlers: make(map[string]HandleFunc, 0)}
}

func (r *Router) addRouter(method, path string, handler HandleFunc) {
	log.Printf("Route %4s - %s", method, path)
	key := method + "," + path
	r.handlers[key] = handler
}

func (r *Router) handle(c *Context) {
	key := c.Req.Method + "," + c.Req.URL.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}

}
