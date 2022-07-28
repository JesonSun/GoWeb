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

	handlers []HandlerFunc
	index    int

	engine *Engine
}

func newContext(req *http.Request, rsp http.ResponseWriter) *Context {
	return &Context{
		Req:    req,
		Rsp:    rsp,
		Method: req.Method,
		Path:   req.URL.Path,
		index:  -1,
	}
}

func (context *Context) Fail(code int, err string) {
	context.index = len(context.handlers)
	context.JSON(code, H{"message": err})
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

func (context *Context) HTML(code int, name string, data interface{}) {
	context.SetHeader("Content-Type", "text/html")
	context.Status(code)
	if err := context.engine.htmlTemplates.ExecuteTemplate(context.Rsp, name, data); err != nil {
		context.Fail(500, err.Error())
	}
}

func (context *Context) Next() {
	context.index++
	for ; context.index < len(context.handlers); context.index++ {
		context.handlers[context.index](context)
	}
}
