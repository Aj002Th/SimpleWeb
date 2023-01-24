package main

import (
	"log"
	"net/http"
	"simpleWeb/sim"
)

func main() {
	engine := sim.New()
	log.Fatal(http.ListenAndServe(":8080", engine))
}
