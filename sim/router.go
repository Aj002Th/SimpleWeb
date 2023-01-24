package sim

import (
	"fmt"
	"strings"
)

type router struct {
	roots    map[string]*node       // root["GET"]、root["POST"]
	handlers map[string]HandlerFunc // handlers["GET-/eho"]
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 路径预处理
// *只能使用在路由末尾
// 其他使用方式不会识别
func parsePath(path string) []string {
	parts := strings.Split(path, "/")
	var ret []string

	for _, part := range parts {
		if part != "" {
			ret = append(ret, part)
			if part[0] == '*' {
				break
			}
		}
	}

	return ret
}

// 路由注册
func (r *router) addRoute(method, path string, handler HandlerFunc) {
	key := method + "-" + path
	parts := parsePath(path)
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.handlers[key] = handler
	r.roots[method].insert(path, parts, 0)
}

// 处理路由
func (r *router) handle(ctx *Context) {
	method := ctx.Req.Method
	path := ctx.Req.URL.Path
	searchParts := parsePath(path)

	n := r.roots[method].search(searchParts, 0)
	if n == nil {
		_, _ = fmt.Fprintf(ctx.Writer, "some error in %s: %s", ctx.Req.URL.Path, "no route found")
		return
	}

	// 获取params参数
	key := method + "-" + n.pattern
	parts := parsePath(n.pattern)
	params := map[string]string{}
	for i, part := range parts {
		if part[0] == ':' && len(part) > 1 {
			params[part[1:]] = searchParts[i]
		}

		if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(searchParts[i:], "/")
			break
		}
	}
	ctx.Params = params

	r.handlers[key](ctx)
}
