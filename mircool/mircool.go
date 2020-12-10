package mircool

import (
	"net/http"
)

type HandleFunc func(c *Context)

type Engine struct {
	*RouterGroup
	router *Router
	groups []*RouterGroup
}

func NewServer() *Engine {
	e := &Engine{router: NewRouter()}
	e.RouterGroup = &RouterGroup{engine: e}
	e.groups = []*RouterGroup{e.RouterGroup}
	return e
}

func (group *RouterGroup) Group(path string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + path,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method, path string, handler HandleFunc) {
	newPath := group.prefix + path
	group.engine.router.addRouter(method, newPath, handler)
}

func (group *RouterGroup) GET(path string, handler HandleFunc) {
	group.addRoute("GET", path, handler)
}

func (group *RouterGroup) POST(path string, handler HandleFunc) {
	group.addRoute("POST", path, handler)
}

func (group *RouterGroup) PUT(path string, handler HandleFunc) {
	group.addRoute("PUT", path, handler)
}

func (group *RouterGroup) DELETE(path string, handler HandleFunc) {
	group.addRoute("DELETE", path, handler)
}

func (group *RouterGroup) PATCH(path string, handler HandleFunc) {
	group.addRoute("PATCH", path, handler)
}

func (e *Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	context := newContext(resp, req)
	e.router.handle(context)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

type RouterGroup struct {
	prefix      string
	middleWares []HandleFunc //中间件支持
	parent      *RouterGroup //父分组
	engine      *Engine
}
