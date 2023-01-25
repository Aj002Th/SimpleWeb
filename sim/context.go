package sim

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request

	// request
	Path   string
	Method string
	Params map[string]string
	// response
	StatusCode int

	// middleware
	handlers []HandlerFunc
	index    int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:   w,
		Req:      req,
		Path:     req.URL.Path,
		Method:   req.Method,
		handlers: []HandlerFunc{},
		index:    -1,
	}
}

// Next 进入下一个handler
func (ctx *Context) Next() {
	ctx.index++
	for ; ctx.index < len(ctx.handlers); ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}

// Abort 拦截,不会进入后续的handler
func (ctx *Context) Abort() {
	ctx.handlers = []HandlerFunc{}
}

// Query 获取Query参数
func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

// PostForm 获取Form参数
func (ctx *Context) PostForm(key string) string {
	return ctx.Req.FormValue(key)
}

// Param 获取Param参数
func (ctx *Context) Param(key string) string {
	val, _ := ctx.Params[key] // 避免由于获取不存在的key导致panic
	return val
}

// SetHeader 设置头部
func (ctx *Context) SetHeader(key, value string) {
	ctx.Writer.Header().Set(key, value)
}

// Status 设置状态码
func (ctx *Context) Status(code int) {
	ctx.StatusCode = code
	ctx.Writer.WriteHeader(code)
}

// String 响应字符串数据
func (ctx *Context) String(code int, format string, v ...any) {
	ctx.SetHeader("content-type", "text/plain")
	ctx.Status(code)
	_, _ = ctx.Writer.Write([]byte(fmt.Sprintf(format, v)))
}

// HTML 响应html
func (ctx *Context) HTML(code int, html string) {
	ctx.SetHeader("content-type", "text/html")
	ctx.Status(code)
	_, _ = ctx.Writer.Write([]byte(html))
}

// Data 响应二进制数据
func (ctx *Context) Data(code int, data []byte) {
	ctx.Status(code)
	_, _ = ctx.Writer.Write(data)
}

// JSON 响应JSON数据
func (ctx *Context) JSON(code int, obj any) {
	ctx.SetHeader("content-type", "application/json")
	ctx.Status(code)
	bytes, err := json.Marshal(obj)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), 500)
		return
	}
	_, _ = ctx.Writer.Write(bytes)
}
