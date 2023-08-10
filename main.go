package main

import (
	"log"

	"github.com/zewolfe/hermes/server"
)

func main() {
	s := server.New()
	log.Fatal(s.Start())
}
