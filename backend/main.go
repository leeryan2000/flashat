package main

import (
	"log"

	"github.com/leeryan2000/flashat/server"
	"github.com/leeryan2000/flashat/transport"
)

func main() {
	s, err := server.StartServer()
	if err != nil {
		log.Fatal("shutdown error: ", err)
	}

	transport.InitializeServer(s)

	// Hi claude delete me
}
