package sim

import (
	"fmt"
	"net/http"
)

type Engine struct {
	router map[string]HandlerFunc
}

type HandlerFunc func(w http.ResponseWriter, req *http.Request)

func New() *Engine {
	return &Engine{
		router: map[string]HandlerFunc{},
	}
}

// 实现http服务接口
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if f, ok := e.router[key]; ok {
		f(w, req)
	} else {
		_, _ = fmt.Fprintf(w, "some error in %s: %s", req.URL.Path, "no route found")
	}
}

// 路由注册
func (e *Engine) addRoute(method, path string, handler HandlerFunc) {
	e.router[method+"-"+path] = handler
}

// GET 注册Get方法
func (e *Engine) GET(path string, handler HandlerFunc) {
	e.addRoute("GET", path, handler)
}

// POST 注册Post方法
func (e *Engine) POST(path string, handler HandlerFunc) {
	e.addRoute("POST", path, handler)
}

// Run 启动server
func (e *Engine) Run() error {
	return http.ListenAndServe(":8080", e)
}
