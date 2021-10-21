package gee

import (
	"net/http"
	"strings"
)

type handleFunc func(c *Context)

type RouterGroup struct {
	prefix      string       //分组前缀
	parent      *RouterGroup //父分组
	engine      *Engine      //任何分组都拥有同一个engine
	middlewares []handleFunc //中间件
}

type Engine struct {
	*RouterGroup // engine拥有第一层分组 类似于 /hello 中的"/" ，继承了routergroup 可以访问routergroup下的方法
	Router       *router
	groups       []*RouterGroup //保存所有分组
}

func New() *Engine {
	engine := &Engine{Router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = append(engine.groups, engine.RouterGroup)
	return engine
}

func (g *RouterGroup) Group(pre string) *RouterGroup {
	group := &RouterGroup{
		engine: g.engine,
		prefix: g.prefix + pre,
		parent: g,
	}
	g.engine.groups = append(g.engine.groups, group)
	return group
}

func (g *RouterGroup) Get(path string, handler handleFunc) {
	g.addRoute("Get", path, handler)
}

func (g *RouterGroup) Use(midlewares ...handleFunc) {
	g.middlewares = append(g.middlewares, midlewares...)
}

func (g *RouterGroup) Post(path string, handler handleFunc) {
	g.addRoute("Post", path, handler)
}

func (engine *Engine) Run(port string) (err error) {
	return http.ListenAndServe(port, engine)
}

//适配器模式？
func (g *RouterGroup) addRoute(method string, pattern string, handler handleFunc) {
	wholePattern := g.prefix + pattern
	g.engine.Router.addRoute(method, wholePattern, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	midlewares := make([]handleFunc, 0)
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			midlewares = append(midlewares, group.middlewares...)
		}
	}

	c := NewContext(w, req)
	c.handlers = midlewares
	engine.Router.handle(c)
}
