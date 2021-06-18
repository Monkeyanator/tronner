package main

import (
	"flag"
	"log"

	"github.com/Monkeyanator/tronner/server"
)

var (
	port uint
)

func main() {
	registerFlags()
	log.Printf("Starting Tron server on port %d", port)
	server := server.New(server.Config{
		Port: port,
	})
	log.Fatal(server.Run())
}

func registerFlags() {
	flag.UintVar(&port, "port", 8080, "port to run server on")
	flag.Parse()
}
