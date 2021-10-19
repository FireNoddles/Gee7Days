package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	W      http.ResponseWriter
	Req    *http.Request
	Code   int
	Path   string
	Method string
	Params map[string]string
}

func (c *Context) Param(key string) string {
	value, ok := c.Params[key]
	if ok {
		return value
	} else {
		return ""
	}
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

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("content-type", "text/plain")
	c.SetStatus(code)
	c.W.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) Json(code int, obj interface{}) {
	c.SetHeader("content-type", "application/json")
	c.SetStatus(code)
	encoder := json.NewEncoder(c.W)
	if err := encoder.Encode(obj); err != nil {
		return
	}
}

func (c *Context) Data(code int, data []byte) {
	c.SetStatus(code)
	c.W.Write(data)
}

func (c *Context) Html(code int, html string) {
	c.SetHeader("content-type", "application/html")
	c.SetStatus(code)
	c.W.Write([]byte(html))
}
