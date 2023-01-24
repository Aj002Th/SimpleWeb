package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"simpleWeb/sim"
)

func main() {
	engine := sim.New()
	engine.POST("/echo", echoHandler)
	log.Fatal(engine.Run())
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	data := string(bytes)
	fmt.Printf("get req: %v", data)
	if err != nil {
		_, _ = fmt.Fprintf(w, "some error in /echo: %v", err)
		if err != nil {
			return
		}
		return
	}
	_, _ = fmt.Fprint(w, data)
}
