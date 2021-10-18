package main

import (
	"./gee"
	"net/http"
)

func main() {
	engine := gee.New()
	engine.Get("/", func(c *gee.Context) {
		c.Html(http.StatusOK, "<h1>gee hello</h1>")
	})
	engine.Get("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello, %s", c.GetValue("name"))
	})
	engine.Post("/login", func(c *gee.Context) {
		c.Json(http.StatusOK,
			gee.H{
				"username": c.GetValue("name"),
				"password": c.GetValue("password"),
			})
	})
	engine.Run(":8999")
}
