package main

import (
	"github.com/id27182/ami-proxy/pkg/server"
	"log"
)

func main() {
	err := server.Serve()
	if err != nil {
		log.Fatal(err)
	}
}