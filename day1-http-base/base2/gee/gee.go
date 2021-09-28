package gee

import (
	"fmt"
	"net/http"
)

type handleFunc func(w http.ResponseWriter, req *http.Request)

type Engine struct {
	Router map[string]handleFunc
}

func New() *Engine{
	return &Engine{Router: make(map[string]handleFunc)}
}

func (engine *Engine) addRouter(method string, path string, handler handleFunc){
	key := method + "-" +path
	engine.Router[key] = handler
}

func (engine *Engine) Get(path string, handler handleFunc){
	engine.addRouter("Get", path, handler)
}

func (engine *Engine) Post(path string, handler handleFunc){
	engine.addRouter("Post", path, handler)
}

func (engine *Engine) Run(port string) (err error){
	return http.ListenAndServe(port, engine)
}


func (engine * Engine) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.Router[key]; ok{
		handler(w, req)
	}else{
		fmt.Fprint(w,"404 not found")
	}
}
