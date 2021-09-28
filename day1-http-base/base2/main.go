package main

import (
	"./gee"
	"fmt"
	"net/http"
)

func main()  {
	engine := gee.New()
	engine.Get("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "URL.Path=%q\n", req.URL.Path)
	})
	engine.Get("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k,v := range req.Header{
			fmt.Fprint(w, "k,v", k, v)
		}
	})
	engine.Run(":8999")
}
