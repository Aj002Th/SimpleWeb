package sim

import (
	"fmt"
	"io"
	"net/http"
)

type Engine struct{}

func New() *Engine {
	return &Engine{}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/echo":
		echoHandler(w, req)
	default:
		_, _ = fmt.Fprintf(w, "some error in %s: %s", req.URL.Path, "no route found")
	}
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
