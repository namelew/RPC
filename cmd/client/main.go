package main

import (
	"github.com/namelew/RPC/internal/rpc"
)

func main() {
	go rpc.Listen()
}
