package main

import (
	"./gee"
	"fmt"
	"net/http"
)

func main()  {
	engine := gee.New()
	engine.Get("/", func(c *gee.Context) {
		c.String()
	})
	engine.Get("/hello", func(c *gee.Context) {
		for k,v := range req.Header{
			fmt.Fprint(w, "k,v", k, v)
		}
	})
	engine.Run(":8999")
}
