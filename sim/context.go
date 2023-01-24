package sim

import "net/http"

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request

	// 常用所以提取出来
	// request
	Path   string
	Method string
	// response
	Code int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}
