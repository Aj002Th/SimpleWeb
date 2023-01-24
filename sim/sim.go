package sim

import (
	"net/http"
)

type HandlerFunc func(ctx *Context)

type H map[string]interface{}

type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

// 实现http服务接口
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := newContext(w, req)
	e.router.handle(ctx)
}

// GET 注册Get方法
func (e *Engine) GET(path string, handler HandlerFunc) {
	e.router.addRoute("GET", path, handler)
}

// POST 注册Post方法
func (e *Engine) POST(path string, handler HandlerFunc) {
	e.router.addRoute("POST", path, handler)
}

// Run 启动server
func (e *Engine) Run() error {
	return http.ListenAndServe(":8080", e)
}
