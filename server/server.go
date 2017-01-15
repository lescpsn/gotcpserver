package main

import (
	"gotcpserver"
)

func main() {
	srv := gotcpserver.NewServer("127.0.0.1:9012")
	srv.Start()
}
