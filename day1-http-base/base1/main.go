package main

import (
	"fmt"
	"net/http"
)
type Engine struct {

}


func (engine *Engine) ServeHTTP (w http.ResponseWriter, req *http.Request){
	switch req.URL.Path {
	case "/":
		fmt.Fprint(w, "URL.Path=%q\n", req.URL.Path)

	}
}