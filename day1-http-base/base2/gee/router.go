package gee

type router struct {
	handlers map[string]handleFunc
}


func NewRouter() *router{
	return &router{handlers: make(map[string]handleFunc)}
}

func (r *router) addRouter(method string, path string, handler handleFunc){
	key := method + "-" +path
	r.handlers[key] = handler
}

func (r *router) handle(c *Context){
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok{
		handler(c)
	}else{
		c.String(404, "404 not found")
	}
}