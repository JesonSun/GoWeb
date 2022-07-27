package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// origin objects
	Req *http.Request
	Rsp http.ResponseWriter
	// request info
	Method string
	Path   string

	Params map[string]string
	//response info
	StatusCode int
}

func newContext(req *http.Request, rsp http.ResponseWriter) *Context {
	return &Context{
		Req:    req,
		Rsp:    rsp,
		Method: req.Method,
		Path:   req.URL.Path,
	}
}

func (context *Context) Param(key string) string {
	value, _ := context.Params[key]
	return value
}

func (context *Context) PostForm(key string) string {
	return context.Req.FormValue(key)
}

func (context *Context) Query(key string) string {
	return context.Req.URL.Query().Get(key)
}

func (context *Context) Status(code int) {
	context.StatusCode = code
	context.Rsp.WriteHeader(code)
}

func (context *Context) SetHeader(key string, value string) {
	context.Rsp.Header().Set(key, value)
}

func (context *Context) String(code int, format string, values ...interface{}) {
	context.SetHeader("Content-Type", "text/plain")
	context.Status(code)
	context.Rsp.Write([]byte(fmt.Sprintf(format, values...)))
}

func (context *Context) JSON(code int, obj interface{}) {
	context.SetHeader("Content-Type", "application/json")
	context.Status(code)
	encoder := json.NewEncoder(context.Rsp)
	if err := encoder.Encode(obj); err != nil {
		http.Error(context.Rsp, err.Error(), 500)
	}
}

func (context *Context) Data(code int, data []byte) {
	context.Status(code)
	context.Rsp.Write(data)
}

func (context *Context) HTML(code int, html string) {
	context.SetHeader("Content-Type", "text/html")
	context.Status(code)
	context.Rsp.Write([]byte(html))
}
