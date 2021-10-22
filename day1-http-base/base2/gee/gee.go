package gee

import (
	"html/template"
	"net/http"
	"path"
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
	*RouterGroup  // engine拥有第一层分组 类似于 /hello 中的"/" ，继承了routergroup 可以访问routergroup下的方法
	Router        *router
	groups        []*RouterGroup     //保存所有分组
	htmlTemplates *template.Template //html渲染
	funcMap       template.FuncMap   //html渲染
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
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
	c.engine = engine
	engine.Router.handle(c)
}

func (g *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) handleFunc {
	abPath := path.Join(g.prefix, relativePath)
	fileServer := http.StripPrefix(abPath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.SetStatus(http.StatusNotFound)
		}
		fileServer.ServeHTTP(c.W, c.Req)
	}
}

func (g *RouterGroup) Static(relativePath string, root string) {
	handler := g.createStaticHandler(relativePath, http.Dir(root))
	//映射文件的路径为relativePath/*filepath 这样文件都会到动态路由relativePath下去找
	urlPattern := path.Join(relativePath, "/*filepath")
	g.Get(urlPattern, handler)
}
