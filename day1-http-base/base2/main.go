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
	case "/hello":
		for k,v := range req.Header{
			fmt.Fprint(w, "k,v", k, v)
		}
	default:
		fmt.Fprint(w, "404")
	}
}

func main()  {
	engine := &Engine{}
	http.ListenAndServe(":8999", engine)
}