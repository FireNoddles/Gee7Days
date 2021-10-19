package gee

import (
	"net/http"
)

type handleFunc func(c *Context)

type Engine struct {
	Router *router
}

func New() *Engine {
	return &Engine{Router: NewRouter()}
}

func (engine *Engine) Get(path string, handler handleFunc) {
	engine.Router.addRoute("Get", path, handler)
}

func (engine *Engine) Post(path string, handler handleFunc) {
	engine.Router.addRoute("Post", path, handler)
}

func (engine *Engine) Run(port string) (err error) {
	return http.ListenAndServe(port, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(w, req)
	engine.Router.handle(c)
}
