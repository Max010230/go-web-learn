package mircool

import (
	"log"
	"net/http"
	"strings"
)

type Router struct {
	roots    map[string]*node
	handlers map[string]HandleFunc
}

func NewRouter() *Router {
	return &Router{
		roots:    make(map[string]*node, 0),
		handlers: make(map[string]HandleFunc, 0),
	}
}

func parsePath(path string) []string {
	split := strings.Split(path, "/")
	parts := make([]string, 0)
	for _, str := range split {
		if str != "" {
			parts = append(parts, str)
			if str[0] == '*' {
				break

			}
		}
	}
	return parts
}

func (r *Router) addRouter(method, path string, handler HandleFunc) {
	parts := parsePath(path)
	log.Printf("Route %4s - %s", method, path)
	key := method + "," + path
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(path, parts, 0)
	r.handlers[key] = handler
}

func (r *Router) getRouter(method, path string) (*node, map[string]string) {
	searchParts := parsePath(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	search := root.search(searchParts, 0)
	if search != nil {
		parts := parsePath(search.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
				break
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return search, params
	}
	return nil, nil
}

func (r *Router) handle(c *Context) {
	router, params := r.getRouter(c.Method, c.Path)
	if router != nil {
		c.Params = params
		key := c.Method + "," + router.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}

}
