package sim

import (
	"net/http"
	"strings"
)

// HandlerFunc sim使用的处理函数
type HandlerFunc func(ctx *Context)

// H 提供便利的简写形式
type H map[string]interface{}

type Engine struct {
	// 把engine抽象成了一个顶级的RouterGroup
	// 路由注册的api直接放入group中
	// 由于使用了匿名嵌套, engine也能使用对应的api
	*RouterGroup

	// 实现路由功能,存储路由信息
	router *router

	// 存储所有group
	groups []*RouterGroup
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Recovery(), Logger())
	return engine
}

// 实现http服务接口
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := newContext(w, req)

	// 将匹配到的中间件加上
	// 这一步会导致性能降低, 因为每个router会走哪些中间件实际上是已经确定了的
	// 但是这里每次会动态找一次
	for _, group := range e.groups {
		if strings.HasPrefix(ctx.Path, group.prefix) {
			ctx.handlers = append(ctx.handlers, group.middleware...)
		}
	}

	e.router.handle(ctx)
}

// Run 启动server
func (e *Engine) Run() error {
	return http.ListenAndServe(":8080", e)
}
