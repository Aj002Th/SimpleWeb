package main

import (
	"fmt"
	"io"
	"log"
	"simpleWeb/sim"
)

func main() {
	engine := sim.New()

	// 回声
	engine.POST("/echo", func(ctx *sim.Context) {
		bytes, err := io.ReadAll(ctx.Req.Body)
		data := string(bytes)
		fmt.Printf("get req: %v", data)
		if err != nil {
			_, _ = fmt.Fprintf(ctx.Writer, "some error in /echo: %v", err)
			if err != nil {
				return
			}
			return
		}
		_, _ = fmt.Fprint(ctx.Writer, data)
	})

	// string的形式返回query参数"data"
	engine.GET("/query", func(ctx *sim.Context) {
		val := ctx.Query("data")
		ctx.String(200, "data:%s", val)
	})

	// json的形式返回form参数"data"
	engine.POST("/form", func(ctx *sim.Context) {
		val := ctx.PostForm("data")
		ctx.JSON(200, sim.H{
			"data": val,
		})
	})

	// 获取html响应
	engine.GET("/html", func(ctx *sim.Context) {
		ctx.HTML(200, "<h1>hello world<h1/>")
	})

	// json形式获取param参数
	engine.GET("/params/:a/:b/c/*d", func(ctx *sim.Context) {
		ctx.JSON(200, ctx.Params)
	})

	// 使用group
	gs := engine.Group("/groups")
	gs.GET("/group1", func(ctx *sim.Context) {
		ctx.JSON(200, sim.H{
			"msg": "visit group1",
		})
	})

	// 在group上使用中间件打印信息
	gs.Use(func(ctx *sim.Context) {
		fmt.Printf("a req is coming: %s\n", ctx.Path)
		ctx.Next()
		fmt.Printf("a req is leaving: %s\n", ctx.Path)
	})

	// 使用Abort拦截
	stop := engine.Group("/stop")
	stop.Use(func(ctx *sim.Context) {
		fmt.Printf("a req is coming and abort: %s\n", ctx.Path)

		ctx.Abort()
		ctx.JSON(200, sim.H{
			"msg": "abort",
		})

		fmt.Printf("a req is leaving and abort: %s\n", ctx.Path)
	})
	stop.GET("/something", func(ctx *sim.Context) {
		ctx.JSON(200, sim.H{
			"msg": "pass the abort",
		})
	})

	log.Fatal(engine.Run())
}
