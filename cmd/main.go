package main

import (
	"log"

	"github.com/leeryan2000/flashat/http"
	"github.com/leeryan2000/flashat/server"
)

func main() {
	s, err := server.StartServer()
	if err != nil {
		log.Fatal("shutdown error: ", err)
	}

	http.InitializeServer(s)
}
