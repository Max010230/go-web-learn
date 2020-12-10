package mircool

import (
	"net/http"
	"path"
	"strings"
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

func (group *RouterGroup) Use(middlewares ...HandleFunc) {
	group.middleWares = append(group.middleWares, middlewares...)
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
	var middlewares []HandleFunc
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middleWares...)
		}
	}
	context := newContext(resp, req)
	context.handlers = middlewares
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

func (group *RouterGroup) createdStaticHandler(relativePath string, fs http.FileSystem) HandleFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Resp, c.Req)
	}
}

func (group *RouterGroup) Static(relativePath, root string) {
	handler := group.createdStaticHandler(relativePath, http.Dir(root))
	urlPatten := path.Join(relativePath, "/*filepath")
	group.GET(urlPatten, handler)
}
