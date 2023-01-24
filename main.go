package main

import (
	"fmt"
	"io"
	"log"
	"simpleWeb/sim"
)

func main() {
	engine := sim.New()
	engine.POST("/echo", echoHandler)
	log.Fatal(engine.Run())
}

func echoHandler(ctx *sim.Context) {
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
}
