package gee

import (
	"net/http"
)

type Context struct {
	W      http.ResponseWriter
	Req    *http.Request
	Code   int
	Path   string
	Method string
}

func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		W:      w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

func (c *Context) GetValue(key string) string {
	if c.Method == "POST" {
		return c.Req.FormValue(key)
	} else if c.Method == "GET" {
		return c.Req.URL.Query().Get(key)
	} else {
		return ""
	}

}

func (c *Context) SetStatus(code int) {
	c.W.WriteHeader(code)
	c.Code = code
}

func (c *Context) SetHeader(key string, value string) {
	c.W.Header().Set(key, value)
}
