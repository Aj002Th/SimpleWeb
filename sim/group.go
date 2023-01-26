package sim

import (
	"net/http"
	"path"
)

type RouterGroup struct {
	prefix string
	engine *Engine // 用来接触其他资源,实现路由注册等工作

	// 记录在整个分组上使用了哪些中间件
	middleware []HandlerFunc
}

// Group 创建一个group
func (g *RouterGroup) Group(prefix string) *RouterGroup {
	group := &RouterGroup{
		prefix: g.prefix + prefix, // 以这种方式实现嵌套路由
		engine: g.engine,
	}

	// 放入全局group池
	g.engine.groups = append(g.engine.groups, group)

	return group
}

// Use 使用中间件
func (g *RouterGroup) Use(middleware ...HandlerFunc) {
	for _, m := range middleware {
		g.middleware = append(g.middleware, m)
	}
}

// 注册路由（追加上前缀再交给router）
func (g *RouterGroup) addRoute(method, path string, handler HandlerFunc) {
	g.engine.router.addRoute(method, g.prefix+path, handler)
}

// GET 注册Get方法
func (g *RouterGroup) GET(path string, handler HandlerFunc) {
	g.addRoute("GET", path, handler)
}

// POST 注册Post方法
func (g *RouterGroup) POST(path string, handler HandlerFunc) {
	g.addRoute("POST", path, handler)
}

// 包装标准库中提供的fileServer来实现静态资源的加载
func (g *RouterGroup) createFileServer(p string, fs http.FileSystem) HandlerFunc {
	urlPrefix := path.Join(g.prefix, p)
	fileServer := http.StripPrefix(urlPrefix, http.FileServer(fs))
	return func(ctx *Context) {
		// 在这里使用context中的filePath参数获取文件路径可以对文件权限进行限制
		//filePath := ctx.Params["filePath"]

		// 使用标准库提供的fileServer功能进行工作
		fileServer.ServeHTTP(ctx.Writer, ctx.Req)
	}
}

// Static 加载静态资源
func (g *RouterGroup) Static(urlPrefix, root string) {
	handler := g.createFileServer(urlPrefix, http.Dir(root))
	url := path.Join(urlPrefix, "/*filePath")
	g.GET(url, handler)
}
