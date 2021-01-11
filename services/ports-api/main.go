package main

import (
	"log"

	"company.com/seaports/services/ports-api/server"
)

func main() {
	server.StartAsync(8080)

	waitChan := make(chan struct{})

	log.Println("Waiting forever")
	<-waitChan
}
