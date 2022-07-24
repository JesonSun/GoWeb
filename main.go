package main

import (
	"GoWeb/gee"
	"fmt"
	"net/http"
)

func main() {
	engine := gee.New()
	engine.GET("/hello", func(req *http.Request, rsp http.ResponseWriter) {
		fmt.Fprintf(rsp, "URL.Path = %q\n", req.URL.Path)
	})
	engine.Run(":9999")
}
