package sim

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

// Logger 简单的日志打印中间件
func Logger() HandlerFunc {
	return func(ctx *Context) {
		fmt.Printf("%s - %s\n", ctx.Method, ctx.Path)
	}
}

// 触发recover时打印堆栈信息
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

// Recovery 防止服务器因panic宕机
func Recovery() HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				ctx.Fail(500, "Internal Server Error")
			}
		}()

		ctx.Next()
	}
}
