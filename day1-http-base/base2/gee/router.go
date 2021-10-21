package gee

import "strings"

type router struct {
	root     map[string]*node //root['GET'],root['POST']....
	handlers map[string]handleFunc
}

func NewRouter() *router {
	return &router{
		root:     make(map[string]*node),
		handlers: make(map[string]handleFunc),
	}
}

//切割路径 遇到*说明后面是通配
func parsePattern(pattern string) []string {
	ps := strings.Split(pattern, "/")
	parts := make([]string, 0)

	for _, part := range ps {
		if part != "" {
			parts = append(parts, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, path string, handler handleFunc) {
	key := method + "-" + path
	parts := parsePattern(path)

	_, ok := r.root[method]
	if !ok {
		r.root[method] = &node{}
	}

	r.root[method].insert(path, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(path string, method string) (*node, map[string]string) {
	parts := parsePattern(path)
	_, ok := r.root[method]
	if !ok {
		return nil, nil
	}
	params := make(map[string]string, 0)
	n := r.root[method].search(parts, 0)

	//拿动态路由的参数params
	if n != nil {
		goalParts := parsePattern(n.pattern)

		for index, part := range goalParts {
			if part[0] == ':' {
				params[part[1:]] = parts[index]
			}

			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(parts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Path, c.Method)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(404, "404 not found")
		})

	}
	c.Next()
}
