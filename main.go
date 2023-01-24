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

	log.Fatal(engine.Run())
}
