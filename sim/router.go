package sim

import "fmt"

type router struct {
	handler map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		handler: map[string]HandlerFunc{},
	}
}

// 路由注册
func (r *router) addRoute(method, path string, handler HandlerFunc) {
	r.handler[method+"-"+path] = handler
}

// 处理路由
func (r *router) handle(ctx *Context) {
	key := ctx.Req.Method + "-" + ctx.Req.URL.Path
	if f, ok := r.handler[key]; ok {
		f(ctx)
	} else {
		_, _ = fmt.Fprintf(ctx.Writer, "some error in %s: %s", ctx.Req.URL.Path, "no route found")
	}
}
