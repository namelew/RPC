package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/namelew/RPC/internal/coordinator"
)

func main() {
	godotenv.Load()
	go coordinator.Listen()
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
}
