package main

import (
	"log"
	"prometheus-transformer/config"
	"prometheus-transformer/server"
)

func main() {
	log.Println("Starting Prometheus Transformer...")
	config.InitConfig()
	server.StartServer()
}
